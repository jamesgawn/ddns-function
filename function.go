package ddnsfunction

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleDDNSUpdate(w http.ResponseWriter, r *http.Request) {
	SendResponse(w, 200, "Hello World!")
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
	buf, err := ioutil.ReadFile("VERSION")
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}
