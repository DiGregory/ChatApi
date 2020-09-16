package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DiGregory/ChatApi/controllers"
	"github.com/DiGregory/ChatApi/storage"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

type apiApp struct {
	addr    string
	Storage *storage.ChatStorage
}

func createApp(addr string, chatStorage *storage.ChatStorage) (*apiApp, error) {
	return &apiApp{addr, chatStorage}, nil
}

func (a apiApp) start() (error) {
	return a.CreateHandlers(a.addr, a.Storage)
}

func (a apiApp) CreateHandlers(addr string, s *storage.ChatStorage) (error) {
	r := chi.NewRouter()
	r.Post("/users/add", func(w http.ResponseWriter, r *http.Request) {
		controllers.AddUser(w, r, s)
	})
	r.Post("/chats/add", func(w http.ResponseWriter, r *http.Request) {
		controllers.AddChat(w, r, s)
	})
	r.Post("/messages/add", func(w http.ResponseWriter, r *http.Request) {
		controllers.AddMessage(w, r, s)
	})
	r.Get("/chats/get", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetChats(w, r, s)
	})
	r.Get("/messages/get", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetMessages(w, r, s)
	})
	fmt.Println("Server started at: ", addr)
	return http.ListenAndServe(addr, r)
}

func main() {
	dbDSN := "host=db user=postgres-dev password=1234 dbname=dev port=5432 sslmode=disable"

	storageConn, err := storage.NewConn("postgres", dbDSN)
	if err != nil {
		log.Fatal("can`t connect to DB: ", err)
	}

	myApp, err := createApp(":9000", storageConn)
	if err != nil {
		log.Fatal("can`t run API application: ", err)
	}
	err = myApp.start()
	if err != nil {
		log.Fatal("can`t start API application: ", err)
	}
}
