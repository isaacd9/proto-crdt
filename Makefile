.PHONY: compile test

compile:
	protoc -I=proto/v1 --go_out=go/pb/v1 --go_opt=module=github.com/isaacd9/proto-crdt/pb/v1 proto/v1/*
	cd go && go build ./...
	cd rust && cargo build --lib

test:
	cd go && go test ./... -v

monotonic_counter:
	protoc -I=proto/v1 -I=examples/monotonic_counter --go-grpc_opt=paths=source_relative --go-grpc_out=examples/monotonic_counter/pb --go_out=examples/monotonic_counter --go_opt=module=github.com/isaacd9/proto-crdt/examples/monotonic_counter examples/monotonic_counter/*.proto

shopping_cart:
	protoc -I=proto/v1 -I=examples/shopping_cart --go-grpc_opt=paths=source_relative --go-grpc_out=examples/shopping_cart/pb --go_out=examples/shopping_cart --go_opt=module=github.com/isaacd9/proto-crdt/examples/shopping_cart examples/shopping_cart/*.proto
