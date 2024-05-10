go-build:
	go build -o ~/bin/kt-crawler cmd/main.go

temporal-dev:
	temporal server start-dev -f /tmp/temporal.db

temporal-worker:
	kt-crawler worker

go-mod-tidy:
	go mod tidy

test:
	go test -v ./...

test.expensive_cases:
	go test -v -tags expensive_tests ./...
