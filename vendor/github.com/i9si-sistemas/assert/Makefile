COVERAGE_DIR="coverage"


test:
	@go test ./... --race -cover -v

cover:
	@mkdir -p $(COVERAGE_DIR) && go test -coverprofile=./$(COVERAGE_DIR)/coverage.out ./...  && go tool cover -html=./$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/index.html