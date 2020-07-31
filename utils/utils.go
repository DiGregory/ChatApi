package utils

import (
	"encoding/json"
	"net/http"
)

func Respond(w http.ResponseWriter, status int, message map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, ok := message["error"]; ok {
		message["status"] = status
	}
	json.NewEncoder(w).Encode(message)

}
