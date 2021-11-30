# this command lets you just run the protoc command to compile proto files
compile-proto:
	@mkdir -p api/doc && docker run --rm -v "$(PWD)/api":/api -w "/api" thethingsindustries/protoc \
		--go_out=. --go_opt=paths=source_relative --go-grpc_out=.  --go-grpc_opt=paths=source_relative --proto_path=. --doc_out=./doc --doc_opt=html,index.html auction.proto

# Start the system
start-system-linux:
	@docker-compose build && docker-compose up

start-system:
	@docker compose build && docker compose up

# Connect a client to the system
start-client:
	@go run client/main.go --address localhost:50000 --name "Bob the Builder"
