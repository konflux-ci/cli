BINARY_NAME=konflux
GIT_COMMIT := $(shell git rev-list -1 HEAD)


build:
	mkdir -p bin
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin -ldflags "-X main.GitCommit=${GIT_COMMIT}" cmd/konflux/main.go
	GOARCH=amd64 GOOS=linux  go build -o bin/${BINARY_NAME}-linux -ldflags "-X main.GitCommit=${GIT_COMMIT}" cmd/konflux/main.go
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows -ldflags "-X main.GitCommit=${GIT_COMMIT}" cmd/konflux/main.go


clean:
	go clean
	rm bin/${BINARY_NAME}-darwin
	rm bin/${BINARY_NAME}-linux
	rm bin/${BINARY_NAME}-windows
