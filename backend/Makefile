EXAMPLE_REQUEST_TALK_PATH = fixtures/request/example_talk.json
EXAMPLE_REQUEST_SUBSCRIBE_PATH = fixtures/request/example_subscribe.json

.PHONY: run
run:
	docker compose up --build

.PHONY: run-backend
run-backend:
	docker compose up --build backend

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

.PHONY: test-grpc
test-grpc:
	buf curl --protocol grpc --http2-prior-knowledge \
	--schema ../proto \
	--data "{\"histories\":[{\"date\":\"2024-12-07\",\"user\":{\"userId\":\"123\",\"name\":\"太郎\",\"role\":2},\"message\":\"今度ご飯にいきませんか？\"}],\"actionKind\":1}" \
	http://localhost:50051/backend.v1.BackendService/Talk

.PHONY: test-grpc-subscribe
test-grpc-subscribe:
	buf curl --protocol grpc --http2-prior-knowledge \
	--schema ../proto \
	--data "{\"platform_type\":1,\"user_id\":\"user_id\"}" \
	http://localhost:50051/backend.v1.BackendService/Subscribe

.PHONY: test-request-talk
test-request-talk:
	buf curl --protocol grpc --http2-prior-knowledge \
		--schema ../proto \
		--data "$$(cat $(EXAMPLE_REQUEST_TALK_PATH))" \
		http://localhost:50051/backend.v1.BackendService/Talk

# bufで動かないとき用
.PHONY: test-request-talk-debug
test-request-talk-debug:
	grpcurl -plaintext \
		-import-path ../proto \
		-proto ../proto/backend/v1/backend.proto \
		-d @ \
		localhost:50051 \
		backend.v1.BackendService/Talk < $(EXAMPLE_REQUEST_TALK_PATH)

.PHONY: test-request-subscribe
test-request-subscribe:
	buf curl --protocol grpc --http2-prior-knowledge \
		--schema ../proto \
		--data "$$(cat $(EXAMPLE_REQUEST_SUBSCRIBE_PATH))" \
		http://localhost:50051/backend.v1.BackendService/Subscribe

# bufで動かないとき用
.PHONY: test-request-subscribe-debug
test-request-subscribe-debug:
	grpcurl -plaintext \
		-import-path ../proto \
		-proto ../proto/backend/v1/backend.proto \
		-d @ \
		localhost:50051 \
		backend.v1.BackendService/Subscribe < $(EXAMPLE_REQUEST_SUBSCRIBE_PATH)

.PHONY: stripe-webhook
stripe-webhook:
	stripe listen --forward-to localhost:50080/stripe/webhook
