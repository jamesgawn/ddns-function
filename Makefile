build:
	mkdir -p dist
	cp VERSION ./dist/VERSION
	go env -w GOFLAGS=-mod=mod
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/main
deploy:
	terraform apply infra -auto-approve