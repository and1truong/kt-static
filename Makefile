temporal-dev:
	temporal server start-dev

go-mod-tidy:
	go mod tidy

test:
	go test -v ./...

test.expensive_cases:
	go test -v -tags expensive_tests ./...
