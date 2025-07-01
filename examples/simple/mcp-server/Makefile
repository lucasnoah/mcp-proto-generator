.PHONY: build run test clean docker

build:
	@go build -o mcp-server .

run: build
	@./mcp-server

test:
	@go test ./...

clean:
	@rm -f mcp-server

docker:
	@docker build -t mcp-server .

docker-run: docker
	@docker run -p 3333:3333 \
		-e MCP_API_KEYS="test-key" \
		mcp-server
