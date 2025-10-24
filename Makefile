pro:
	gemini --model gemini-2.5-pro

flash:
	gemini --model gemini-flash-latest

go-test:
	go test -race -v -timeout 12s ./...

go-vest:
	go vet ./...

go-fmt:
	go fmt ./...

go-lint:
	golangci-lint run ./...

go-run:
	go run ./cmd/kt-crawler -config config.sample.json
