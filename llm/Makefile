.PHONY: run
run:
	docker compose up --build

.PHONY: test-grpc
test-grpc:
	buf curl --protocol grpc --http2-prior-knowledge \
	--schema ../proto \
	--data "{\"histories\":[ \
		{\"date\":\"2025-01-12T15:00:00+0900\",\"user\":{\"userId\":\"001\",\"name\":\"山田 太郎\",\"role\":1},\"message\":\"こんにちは！\"}, \
		{\"date\":\"2025-01-12T15:01:00+0900\",\"user\":{\"userId\":\"002\",\"name\":\"佐藤 花子\",\"role\":2},\"message\":\"こんにちは、今日はいい天気ですね！\"}, \
		{\"date\":\"2025-01-12T15:02:00+0900\",\"user\":{\"userId\":\"001\",\"name\":\"山田 太郎\",\"role\":1},\"message\":\"そうですね、散歩でもしたくなりますね。\"}, \
		{\"date\":\"2025-01-12T15:03:00+0900\",\"user\":{\"userId\":\"002\",\"name\":\"佐藤 花子\",\"role\":2},\"message\":\"近くの公園に新しいカフェができたみたいですよ。\"}, \
		{\"date\":\"2025-01-12T15:04:00+0900\",\"user\":{\"userId\":\"001\",\"name\":\"山田 太郎\",\"role\":1},\"message\":\"本当ですか？気になりますね。どんなカフェか知っていますか？\"}, \
		{\"date\":\"2025-01-12T15:05:00+0900\",\"user\":{\"userId\":\"002\",\"name\":\"佐藤 花子\",\"role\":2},\"message\":\"雑貨と本が楽しめるカフェみたいです。一緒に行きませんか？\"} \
	],\"actionKind\":2}" \
	http://localhost:50052/llm.v1.LlmService/Talk
