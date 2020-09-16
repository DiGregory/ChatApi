package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"

	"github.com/lib/pq"
)

type ChatStorage struct {
	DB *sql.DB
}

func NewConn(driver, source string) (*ChatStorage, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	initQuery, err := ioutil.ReadFile("./storage/init.sql")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(string(initQuery))
	if err != nil {
		return nil, err
	}

	fmt.Println("Connection with db was set up")
	return &ChatStorage{DB: db}, nil
}

func (s *ChatStorage) AddUser(rawUser []byte) (*User, error) {
	var u User
	err := json.Unmarshal(rawUser, &u)
	if err != nil {
		return nil, err
	}

	err = u.validate()
	if err != nil {
		return nil, err
	}

	userRow := s.DB.QueryRow("INSERT INTO users (id,username) VALUES (DEFAULT,$1) RETURNING id;", u.Username)

	var user User
	err = userRow.Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u User) validate() (error) {
	if len(u.Username) == 0 {
		return errors.New("invalid username")
	}
	return nil
}

func (s *ChatStorage) AddChat(rawChat []byte) (*Chat, error) {
	var c Chat
	err := json.Unmarshal(rawChat, &c)
	if err != nil {
		return nil, err
	}

	err = c.validate(s)
	if err != nil {
		return nil, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	chatRow := tx.QueryRow("INSERT INTO chats (id,name) VALUES (DEFAULT,$1) RETURNING id;", c.Name)
	var chat Chat
	err = chatRow.Scan(&chat.ID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, v := range c.Users {
		_, err := tx.Exec("INSERT INTO user_chats(chat_id,user_id) VALUES ($1,$2);", chat.ID, v)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()

	return &chat, nil
}

func (c Chat) validate(s *ChatStorage) (error) {
	if len(c.Name) == 0 {
		return errors.New("invalid chat name")
	}
	if len(c.Users) < 2 {
		return errors.New("invalid number of users")
	}
	//check existing
	for _, v := range c.Users {
		userRow := s.DB.QueryRow("SELECT (id) FROM users WHERE id=$1;", v)
		var user User
		err := userRow.Scan(&user.ID)
		if err != nil {
			return errors.New(fmt.Sprintf("there is no user with id=%v", v))
		}
	}
	return nil
}

func (s *ChatStorage) AddMessage(rawMessage []byte) (*Message, error) {
	var m Message
	err := json.Unmarshal(rawMessage, &m)
	if err != nil {
		return nil, err
	}

	err = m.validate(s)
	if err != nil {
		return nil, err
	}

	messageRow := s.DB.QueryRow("INSERT INTO messages (id,chat_id,user_id,text) VALUES (DEFAULT,$1,$2,$3) RETURNING id;", m.ChatID, m.UserID, m.Text)

	var message Message
	err = messageRow.Scan(&message.ID)
	if err != nil {
		return nil, err
	}

	return &message, nil

}

func (m Message) validate(s *ChatStorage) (error) {
	//check existing
	userRow := s.DB.QueryRow("SELECT (id) FROM users WHERE id=$1;", m.UserID)
	var user User
	err := userRow.Scan(&user.ID)
	if err != nil {
		return errors.New(fmt.Sprintf("there is no user with id=%v", m.UserID))
	}

	chatRow := s.DB.QueryRow("SELECT (id) FROM chats WHERE id=$1;", m.ChatID)
	var chat Chat
	err = chatRow.Scan(&chat.ID)
	if err != nil {
		return errors.New(fmt.Sprintf("there is no chat with id=%v", m.ChatID))
	}

	//does the user exist in the chat
	userRow = s.DB.QueryRow("SELECT (user_id) FROM user_chats WHERE chat_id=$1 AND user_id=$2;", m.ChatID, m.UserID)
	var u User
	err = userRow.Scan(&u.ID)
	if err != nil {
		return errors.New("user is not in chat")
	}
	return nil
}

func (s *ChatStorage) GetChats(rawUser []byte) ([]Chat, error) {
	var User map[string]interface{}
	err := json.Unmarshal(rawUser, &User)
	if err != nil {
		return nil, err
	}
	id, ok := User["user"].(string)
	if !ok {
		return nil, errors.New("wrong input user")
	}
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	//get chats id
	userRows, err := s.DB.Query("SELECT * FROM user_chats WHERE user_id=$1;", uid)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	var userChat UserChat
	userChats := make([]UserChat, 0)
	for userRows.Next() {
		err := userRows.Scan(&userChat.ChatID, &userChat.UserID)
		if err != nil {
			return nil, err
		}
		userChats = append(userChats, userChat)
	}

	var ids []int
	for _, v := range userChats {

		cid, err := strconv.Atoi(v.ChatID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, cid)
	}

	chats := make([]Chat, 0)
	chatRows, err := s.DB.Query("SELECT * FROM chats WHERE id = ANY($1) ORDER BY (SELECT created_at FROM messages WHERE chat_id=chats.id ORDER by created_at DESC LIMIT 1) DESC;", pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer chatRows.Close()
	for chatRows.Next() {
		var c Chat
		err := chatRows.Scan(&c.ID, &c.Name, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		chats = append(chats, c)
	}

	//get users id
	for i, v := range chats {
		chatRows, err := s.DB.Query("SELECT (user_id) FROM user_chats WHERE chat_id=$1;", v.ID)
		if err != nil {
			return nil, err
		}
		defer chatRows.Close()
		userIDs := make([]string, 0)
		for chatRows.Next() {
			d := 0
			err := chatRows.Scan(&d)
			if err != nil {
				return nil, err
			}
			userIDs = append(userIDs, strconv.Itoa(d))
		}
		chats[i].Users = append(chats[i].Users, userIDs...)
		chatRows.Close()
	}
	return chats, nil
}

func (s *ChatStorage) GetMessages(rawChat []byte) ([]Message, error) {
	var msg Message
	err := json.Unmarshal(rawChat, &msg)
	fmt.Println(msg)
	if err != nil {
		return nil, err
	}
	messagesRows, err := s.DB.Query("SELECT * FROM messages WHERE chat_id=$1;", msg.ChatID)
	if err != nil {
		return nil, err
	}
	defer messagesRows.Close()

	var m Message
	messages := make([]Message, 0)
	for messagesRows.Next() {
		err := messagesRows.Scan(&m.ID, &m.ChatID, &m.UserID, &m.Text, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		fmt.Println(m)
		messages = append(messages, m)
	}
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.Before(*messages[j].CreatedAt)
	})
	return messages, nil

}
