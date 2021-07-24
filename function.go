package ddnsfunction

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/" {
		HandleRootPath(w)
	} else if r.Method == http.MethodGet && r.URL.Path == "/nic/update" {
		HandleDDNSUpdate(w, r)
	} else {
		SendResponse(w, 404, "Not Found")
	}
}

func HandleRootPath(w http.ResponseWriter) {
	version := ObtainVersion()
	body := fmt.Sprintf("Dynamic DNS Service (%s)", version)
	SendResponse(w, 200, body)
}

func HandleDDNSUpdate(w http.ResponseWriter, r *http.Request) {
	providedUsername, providedPassword, ok := r.BasicAuth()
	if !ok {
		SendResponse(w, 401, "badauth")
	}

	_, ok = r.URL.Query()["myip"]
	if !ok {
		SendResponse(w, 400, "No IP address provided")
	}

	_, ok = r.URL.Query()["hostname"]
	if !ok {
		SendResponse(w, 400, "No hostname provided")
	}

	username := os.Getenv("username")
	password := os.Getenv("password")
	if username != providedUsername || password != providedPassword {
		SendResponse(w, 401, "badauth")
	}

	SendResponse(w, 501, "Not Implemented")
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
	if statusCode >= 200 && statusCode < 400 {
		os.Exit(0)
	} else {
		os.Exit(1)
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
