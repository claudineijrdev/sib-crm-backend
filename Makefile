run:
	go run cmd/server/main.go

build:
	go build -o crm cmd/server/main.go

test:
	go test ./...
