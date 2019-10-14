package utils

import (
	"encoding/json"
	"net/http"
)

//Message builds a JSON response when called
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

//Respond to the frontend
func Respond(responseWriter http.ResponseWriter, data map[string]interface{}) {
	responseWriter.Header().Add("content-type", "application/json")
	json.NewEncoder(responseWriter).Encode(data)
}
