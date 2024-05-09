temporal-dev:
	temporal server start-dev

go-mod-tidy:
	go mod tidy

test:
	go test -v ./...
