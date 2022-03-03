dependencies:
	go mod vendor
	go mod download
	
build:
	go build -o bin/imgur-go

test:
	go test ./... -v -coverprofile cp.out
	go get github.com/t-yuki/gocover-cobertura
	go tool cover -html=cp.out -o cp.html && gocover-cobertura < cp.out > cp.xml