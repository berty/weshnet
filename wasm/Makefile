GO?=go1.19

.PHONY: serve
serve: weshnet.wasm
	$(GO) run ./serve

.PHONY: weshnet.wasm
weshnet.wasm:
	GOOS=js GOARCH=wasm $(GO) build -o weshnet.wasm -tags nofuse .
