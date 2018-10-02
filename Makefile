main:
	dep ensure
	go build -o bin/imgur-go

test:
	go test ./... -v -coverprofile cp.out