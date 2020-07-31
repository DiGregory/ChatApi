package storage

import "time"

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	CreatedAt *time.Time `json:"-"`
}
type Chat struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Users     []string   `json:"users,string"`
	CreatedAt *time.Time `json:"-"`
}
type Message struct {
	ID        int
	ChatID    int
	UserID    int
	Text      string
	CreatedAt *time.Time
}
