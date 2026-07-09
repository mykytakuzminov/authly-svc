package http

import (
	"encoding/json"
	"log"
	"net/http"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, reqBody any) error {
	return json.NewDecoder(r.Body).Decode(reqBody)
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}
