package ddnsfunction

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
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
	ip, ipOk := r.URL.Query()["myip"]
	hostname, hostnameOk := r.URL.Query()["hostname"]
	username, usernameErr := GetUsername()
	password, passwordErr := GetPassword()

	if !basicOk {
		SendResponse(w, 401, "badauth")
	} else if !ipOk {
		SendResponse(w, 400, "No IP address provided")
	} else if !hostnameOk {
		SendResponse(w, 400, "No hostname provided")
	} else if usernameErr != nil {
		SendResponse(w, 500, "")
		log.Fatal(usernameErr)
	} else if passwordErr != nil {
		SendResponse(w, 500, "")
		log.Fatal(passwordErr)
	} else if username != providedUsername || password != providedPassword {
		SendResponse(w, 401, "badauth")
	} else {
		log.Printf("Starting update for %s to %s", hostname, ip)
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

func GetUsername() (string, error) {
	return GetSecret(os.Getenv("USERNAME_SECRET"))
}

func GetPassword() (string, error) {
	return GetSecret(os.Getenv("PASSWORD_SECRET"))
}

func GetSecret(name string) (string, error) {
	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}

	defer func(c *secretmanager.Client) {
		err := c.Close()
		if err != nil {
			log.Fatal("Unable to close connection to secret manager.")
		}
	}(c)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := c.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", err
	}

	return string(result.Payload.Data), nil
}

func ObtainVersion() string {
	version := os.Getenv("VERSION")
	if version == "" {
		return "0.0.0"
	} else {
		return version
	}
}
