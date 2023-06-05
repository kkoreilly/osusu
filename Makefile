.DEFAULT_GOAL := run

build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/client
	go build -o ./bin/server ./cmd/server

run: build
	./bin/server

buildscraper:
	go build -o ./bin/scraper ./cmd/scraper 

runscraper: buildscraper
	./bin/scraper

scrape: runscraper

buildclassify:
	go build -o ./bin/classify ./cmd/classify

runclassify: buildclassify
	./bin/classify

classify: runclassify