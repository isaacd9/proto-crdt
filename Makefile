.PHONY: compile test

compile:
	protoc -I=proto/v1 --go_out=go/pb/v1 --go_opt=module=github.com/isaacd9/proto-crdt/pb/v1 proto/v1/*
	cd go && go build ./...

test:
	cd go && go test ./... -v
