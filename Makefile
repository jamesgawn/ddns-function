build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/main
deploy:
	terraform apply -auto-approve