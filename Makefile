main:
	dep ensure
	go build -o bin/imgur-go

test:
	go test -cover ./... -v