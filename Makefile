.PHONY: build run-example clean help

help:
	@echo "Ethereum Calldata Analysis Tool"
	@echo ""
	@echo "Commands:"
	@echo "  make build          - Build the tool"
	@echo "  make run-example    - Run example analysis (requires Etherscan API key)"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "Usage:"
	@echo "  ./eth-analysis -config config.json -output text"
	@echo "  ./eth-analysis -config config.yml -output csv"
	@echo ""

build:
	@echo "Building eth-analysis..."
	go mod tidy
	go build -o eth-analysis
	@echo "✓ Build complete: ./eth-analysis"

run-example:
	@echo "Running example analysis..."
	@if [ -f config.json ]; then \
		./eth-analysis -config config.json -output text; \
	else \
		echo "✗ config.json not found. Copy config.example.json to config.json first and add your Etherscan API key."; \
	fi

clean:
	rm -f eth-analysis
	go clean
	@echo "✓ Clean complete"

.PHONY: test
test:
	go test -v ./...
