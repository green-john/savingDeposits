package transport

import (
	"encoding/json"
	"log"
	"net/http"
)

// Utility function to respond to http requests.
func respond(w http.ResponseWriter, status int, data interface{}) {
	if p, ok := data.(Public); ok {
		data = p.Public()
	}

	parsedJson, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, err = w.Write(parsedJson)
	if err != nil {
		log.Println("Error responding:", err)
	}
}
