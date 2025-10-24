pro:
	gemini --model gemini-2.5-pro

flash:
	gemini --model gemini-flash-latest

claude:
	ANTHROPIC_MODEL="gpt-oss-20b" && claude

go-test:
	go test -race -v -timeout 12s ./...

go-vest:
	go vet ./...

go-fmt:
	go fmt ./...

go-lint:
	golangci-lint run ./...

go-run-read:
	go run ./cmd/kt-crawler -config config.sample.json read

go-run-scan:
	go run ./cmd/kt-crawler -config config.sample.json scan

go-run-help:
	go run ./cmd/kt-crawler -config config.sample.json -h
