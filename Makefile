.PHONY=run
run: tidy
	@set -a && . .env && set +a && go run ./cmd/bot

.PHONY=lint
lint: tidy
	@.bin/golangci-lint run --config=.golangci-lint.yml ./...

.PHONY=lint-fix
lint-fix: tidy
	@.bin/golangci-lint run --config=.golangci-lint.yml --fix ./...

.PHONY=test
test: tidy
	@go test ./...

.PHONY=install-lint
install-lint:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b .bin v2.2.2

.PHONY=tidy
tidy:
	@go mod tidy