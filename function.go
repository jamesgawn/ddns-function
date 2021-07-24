package ddnsfunction

import (
	"fmt"
	"log"
	"net/http"
)

func HandleDDNSUpdate(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, 200, "Hello World!")
}

func sendResponse(w http.ResponseWriter, statusCode int, body string) {
	w.WriteHeader(statusCode)
	_, err := fmt.Fprint(w, body)
	if err != nil {
		log.Fatal(err)
	}
}
