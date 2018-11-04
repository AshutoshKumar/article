package util

import (
	"encoding/json"
	"net/http"
)

//RespondWithJson wraps a message into json and returns it in the response,along with a header and response code
func RespondWithJson(w http.ResponseWriter, code int, cargo interface{}) {
	response, _ := json.Marshal(cargo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//RespondWithError uses RespondWithJson to handle bad or unserveable requests
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJson(w, code, map[string]string{"error": message})
}

func Error(w http.ResponseWriter, code int, message string) {
	Json(w, code, map[string]string{"error": message})
}

func Json(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func JsonWithCookie(w http.ResponseWriter, code int, payload interface{}, cookie http.Cookie) {
	response, _ := json.Marshal(payload)
	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
