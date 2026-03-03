.PHONY: build wasm test

build:
	go build -o bananascript ./cmd/bananascript

wasm:
	GOOS=js GOARCH=wasm go build -o bananascript.wasm ./cmd/wasm
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" .

test:
	go test ./...
