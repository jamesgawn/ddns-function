package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	route53helper "github.com/jamesgawn/route53-helper"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	lambda.Start(handler)
}

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.WithField("request", req).Info("Starting to handle request")

	switch req.RouteKey {
	case "GET /":
		return handleRootPath(), nil
	case "GET /nic/update":
		return handleDDNSUpdate(req), nil
	default:
		return buildResponse(404, "Not Found"), nil
	}
}

func handleRootPath() events.APIGatewayV2HTTPResponse {
	version := obtainVersion()
	body := fmt.Sprintf("Dynamic DNS Service (%s)", version)
	return buildResponse(200, body)
}

func handleDDNSUpdate(req events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	authenticated := authenticate(req.Headers["authorization"])

	ip := req.QueryStringParameters["myip"]
	hostname := req.QueryStringParameters["hostname"]
	if !authenticated {
		return buildResponse(401, "badauth")
	} else if ip == "" {
		return buildResponse(400, "No IP address provided")
	} else if hostname == "" {
		return buildResponse(400, "No hostname provided")
	} else {
		log.Printf("Starting update for %s to %s", hostname, ip)

		_, err := route53helper.GetClient()
		if err != nil {
			log.Error(err)
			return buildResponse(500, "Configuration Error")
		}

		return buildResponse(501, "Not Implemented")
	}
}

func authenticate(authHeader string) bool {
	if authHeader != "" {
		splitAuthValue := strings.SplitN(authHeader, " ", 2)
		if len(splitAuthValue) != 2 || splitAuthValue[0] != "Basic" {
			log.Warnln("Unable to find basic label or value from authorization header")
			return false
		}
		authHeaderDecoded, authHeaderErr := base64.StdEncoding.DecodeString(splitAuthValue[1])
		authHeader = string(authHeaderDecoded)
		authSplit := strings.Split(string(authHeader), ":")
		username := os.Getenv("username")
		password := os.Getenv("password")
		providedUsername := authSplit[0]
		providedPassword := authSplit[1]

		if authHeaderErr != nil {
			log.Error(authHeaderErr)
			log.Warnln("Unable to decode authorization header")
			return false
		} else if username == providedUsername && password == providedPassword {
			log.Infoln("Passed authorization")
			return true
		} else {
			log.WithFields(log.Fields{
				"username":         username,
				"password":         password,
				"providedUsername": providedUsername,
				"providedPassword": providedPassword,
			}).Info("Creds check")
			log.Infoln("Incorrect user/pass provided")
			return false
		}
	} else {
		log.Warnln("No authorization header")
		return false
	}

}

func buildResponse(statusCode int, body string) events.APIGatewayV2HTTPResponse {
	version := obtainVersion()
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers: map[string]string{
			"Version": version,
		},
	}
}

func obtainVersion() string {
	versionRaw, err := ioutil.ReadFile("VERSION")
	version := string(versionRaw)

	if err != nil || version == "" {
		return "0.0.0"
	} else {
		return string(version)
	}
}
