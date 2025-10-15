# protobuf specifications

upon modifications, to regenerate the code, run (**directory where the command is ran from matters btw**):
```sh
$ cd $PROJECT_ROOT/proto
$ export FILES="store/store.proto inventory/inventory.proto armory/armory.proto chat/chat.proto repo/repo.proto"
$ protoc -I. -I$(dirname "$(which protoc)")/../include --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. store/store.proto inventory/inventory.proto armory/armory.proto chat/chat.proto repo/repo.proto
```

`-I$(dirname "$(which protoc)")/../include` is necessary to import `google.protobuf` packages. requires some tweaks if `protoc` is installed via `winget` (on windows):
- navigate to `$(dirname "$(which protoc)")`
- open folder location where protoc.symlink points
- add protoc `/bin` directory to `PATH`
- remove the protoc symlink
- restart shell/terminal (or use choco's `refreshenv`, whatever)


note:
vscode intellisense for protobuf works weirdly for imports, requires further configuration to work properly, so if it shows any error underlines, ignore them. they'll work just fine during compile