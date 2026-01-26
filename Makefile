.PHONY: tailwindcss
VERSION := $(shell git describe --tags --abbrev=0)
APP := vis

chart.js:
	curl https://cdnjs.cloudflare.com/ajax/libs/Chart.js/4.5.0/chart.umd.js > chart.js

build: chart.js
	go build .

release: chart.js
	GOOS=linux GOARCH=amd64 \
		go build -tags release -ldflags "-X main.version=$(VERSION)" -o $(APP)
