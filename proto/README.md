# protobuf specifications

to regenerate the code, run:
```sh
$ set TARGET_FILE=chat.proto
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $TARGET_FILE
```