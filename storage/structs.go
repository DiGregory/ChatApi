package storage

import "time"

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	CreatedAt *time.Time `json:"created_at"`
}
type Chat struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Users     []string   `json:"users,string"`
	CreatedAt *time.Time `json:"created_at"`
}
type Message struct {
	ID        int        `json:"id"`
	ChatID    int        `json:"chat,string"`
	UserID    int        `json:"author,string"`
	Text      string     `json:"text"`
	CreatedAt *time.Time `json:"created_at"`
}
type UserChat struct {
	ChatID string
	UserID string
}
