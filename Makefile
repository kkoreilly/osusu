.DEFAULT_GOAL := run

build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build -o ./osusu

run: build
	./osusu