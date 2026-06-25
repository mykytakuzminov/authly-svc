# ── Code quality ──────────────────────────────────────────
check:
	go mod tidy
	gofmt -w .
	go vet ./...
	golangci-lint run
	go build ./...
