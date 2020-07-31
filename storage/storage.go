package storage

import (
	"database/sql"
	"os"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"errors"
)

type ChatStorage struct {
	DB *sql.DB
}

func NewConn(driver, source string) (*ChatStorage, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	initSQLFile, err := os.Open("./storage/init.sql")
	if err != nil {
		return nil, err
	}
	defer initSQLFile.Close()

	initQuery, err := ioutil.ReadAll(initSQLFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Query(string(initQuery))
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

func (s ChatStorage) AddChat(rawChat []byte) (*Chat, error) {
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
		return nil,err
	}
	chatRow := tx.QueryRow("INSERT INTO chats (id,name) VALUES (DEFAULT,$1) RETURNING id;", c.Name)
	var chat Chat
	err = chatRow.Scan(&chat.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, v := range c.Users {
		rows,err:=tx.Query("INSERT INTO user_chats(chat_id,user_id) VALUES ($1,$2);", c.ID, v)
		rows.Close()
		if err!=nil{
			fmt.Println("hui")
			tx.Rollback()
			return nil,err
		}
	}
	tx.Commit()



	return &chat, nil
}

func (c Chat) validate(s ChatStorage) (error) {
	if len(c.Name) == 0 {
		return errors.New("invalid chat name")
	}
	if len(c.Name) < 2 {
		return errors.New("not enough users")
	}
	//check existing
	for _, v := range c.Users {
		userRow := s.DB.QueryRow("SELECT (username) FROM users WHERE id=$1;", v)
		var user User
		err := userRow.Scan(&user.Username)
		if err != nil {
			return err
		}
		if user.Username == "" {
			return errors.New("there is no user with id=" + v)
		}
	}

	return nil
}
