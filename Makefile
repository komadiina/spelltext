generate:
	mkdir -p ./proto/v1 && \
	protoc --go_out=./proto/v1 --go_opt=module=v1 ./proto/*.proto
