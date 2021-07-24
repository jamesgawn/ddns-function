package ddnsfunction

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func HandleDDNSUpdate(w http.ResponseWriter, r *http.Request) {
	version := ObtainVersion()
	if r.Method == http.MethodGet && r.URL.Path == "/" {
		body := fmt.Sprintf("Dynamic DNS Service (%s)", version)
		SendResponse(w, 200, body)
	} else {
		SendResponse(w, 404, "Not Found")
	}
}

func SendResponse(w http.ResponseWriter, statusCode int, body string) {
	version := ObtainVersion()
	w.Header().Set("Version", version)
	w.WriteHeader(statusCode)
	w.Header().Set("Version", "0.0.0.0")
	_, err := fmt.Fprint(w, body)
	if err != nil {
		log.Fatal(err)
	}
}

func ObtainVersion() string {
	version := os.Getenv("VERSION")
	if version == "" {
		return "0.0.0"
	} else {
		return version
	}
}
