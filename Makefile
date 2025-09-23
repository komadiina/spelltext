generate:
	mkdir -p ./proto/v1 && \
	protoc \
		--proto_path=./proto \
		--go_out=./proto --go-grpc_out=./proto \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		proto/v1/gameserver/gameserver.proto
