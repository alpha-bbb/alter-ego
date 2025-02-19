.PHONY: run
run:
	docker compose up --build

.PHONY: test-grpc
test-grpc:
	buf curl --protocol grpc --http2-prior-knowledge \
	--schema ../proto \
	--data "{\"histories\":[{\"date\":\"2024-12-07\",\"user\":{\"userId\":\"123\",\"name\":\"太郎\",\"role\":2},\"message\":\"今度ご飯にいきませんか？\"}],\"actionKind\":2}" \
	http://localhost:50051/backend.v1.BackendService/Talk

.PHONY: tools
tools:
	@if ! which golangci-lint > /dev/null; then \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.1; \
	else \
		echo "golangci-lint already installed, skipping..."; \
	fi

.PHONY: lint
lint: tools
	golangci-lint run ./... --timeout=5m

.PHONY: format
format: tools
	golangci-lint run ./... --fix
