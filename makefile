build:
	GOOS=darwin GOARCH=amd64 go build -o send
test:
	go test -v ./...