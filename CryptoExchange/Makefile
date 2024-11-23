TEST_DIR="tests"

run:
	@go run cmd/CryptoExchange/main.go --config config/config.yaml
build:
	@go build cmd/CryptoExchange/main.go -o CryptoExchange
coverage-html:
	@cd $(TEST_DIR) && go test -coverprofile=coverage.out -coverpkg=../... -timeout 30s -run ^Test app/*_test.go
	@cd $(TEST_DIR) && sed -i '/suite\|app\|token\|response\|config\|recoverer\|logger/d' coverage.out
	@go tool cover -html=$(TEST_DIR)/coverage.out -o $(TEST_DIR)/coverage.html
	@rm $(TEST_DIR)/coverage.out