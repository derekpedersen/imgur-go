dependencies:
	go mod vendor
	go mod download
	
build:
	go build -o bin/imgur-go

test:
	make test-style
	go test ./... -v -coverprofile cp.out
	go tool cover -html=cp.out -o cp.html && go run github.com/t-yuki/gocover-cobertura@latest < cp.out > cp.xml

test-style:
	bash ./scripts/check-table-tests.sh