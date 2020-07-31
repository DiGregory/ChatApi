package controllers

import (
	"net/http"
	"github.com/DiGregory/avitoTech/storage"
	"io/ioutil"
	u "github.com/DiGregory/avitoTech/utils"
)

func AddUser(w http.ResponseWriter, r *http.Request, s *storage.ChatStorage) {
	var response=make(map[string]interface{})

	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response["error"]=err.Error()
		u.Respond(w, http.StatusBadRequest, response)
		return
	}

	responseData, err := s.AddUser(request)
	if err != nil {
		response["error"]=err.Error()
		u.Respond(w, http.StatusInternalServerError, response)
		return
	}
	response["id"]=responseData.ID
	u.Respond(w, http.StatusCreated, response)
}

func AddChat(w http.ResponseWriter, r *http.Request, s *storage.ChatStorage){
	var response=make(map[string]interface{})

	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response["error"]=err.Error()
		u.Respond(w, http.StatusBadRequest, response)
		return
	}

	responseData, err := s.AddChat(request)
	if err != nil {
		response["error"]=err.Error()
		u.Respond(w, http.StatusInternalServerError, response)
		return
	}
	response["id"]=responseData.ID
	u.Respond(w, http.StatusCreated, response)

}