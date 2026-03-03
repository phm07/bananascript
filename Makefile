.PHONY: build wasm test

build:
	go build -o bananascript ./cmd/bananascript

wasm:
	GOOS=js GOARCH=wasm go build -o bananascript.wasm ./cmd/wasm

test:
	go test ./...
