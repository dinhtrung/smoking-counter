# note: call scripts from /scripts
#GOBUILD=GOOS=linux GOARCH=amd64 go build
#GOWIN=GOOS=windows GOARCH=amd64 go build
GOBUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"'


# Clean all binaries
clean:
	rm -f build/apps/*
	rm -f build/tools/*
	rm -f build/utils/*

go-update:
	go get -u ./...


# Microservice API Gateway
smoking-counter-api:
	$(GOBUILD) -o build/smoke-counter-api cmd/smoke-counter-api/*.go

# Account creation tools
account-cli:
	$(GOBUILD) -o build/account-cli tools/account-cli/*.go
