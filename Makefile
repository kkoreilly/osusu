.DEFAULT_GOAL := run

subset=false
lr=0.01
layers=1
units=500

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
	./bin/classify -subset=$(subset) -lr=$(lr) -layers=$(layers) -units=$(units)

classify: runclassify