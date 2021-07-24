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
	providedUsername, providedPassword, basicOk := r.BasicAuth()
	_, ipOk := r.URL.Query()["myip"]
	_, hostnameOk := r.URL.Query()["hostname"]
	username := GetUsername()
	password := GetPassword()

	if !basicOk {
		SendResponse(w, 401, "badauth")
	} else if !ipOk {
		SendResponse(w, 400, "No IP address provided")
	} else if !hostnameOk {
		SendResponse(w, 400, "No hostname provided")
	} else if username != providedUsername || password != providedPassword {
		SendResponse(w, 401, "badauth")
	} else {
		SendResponse(w, 501, "Not Implemented")
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

func GetUsername() string {
	return os.Getenv("username")
}

func GetPassword() string {
	return os.Getenv("password")
}

func ObtainVersion() string {
	version := os.Getenv("VERSION")
	if version == "" {
		return "0.0.0"
	} else {
		return version
	}
}
