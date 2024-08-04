.PHONY: build test

build:
	go build ./...

test: bin/sqlc-gen-perl.wasm
	go test ./...

all: bin/sqlc-gen-perl bin/sqlc-gen-perl.wasm

bin/sqlc-gen-perl: bin go.mod go.sum $(wildcard **/*.go)
	cd plugin && go build -o ../bin/sqlc-gen-perl ./main.go

bin/sqlc-gen-perl.wasm: bin/sqlc-gen-perl
	cd plugin && GOOS=wasip1 GOARCH=wasm go build -o ../bin/sqlc-gen-perl.wasm main.go

bin:
	mkdir -p bin
