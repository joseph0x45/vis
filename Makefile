.PHONY: tailwindcss
VERSION := $(shell git describe --tags --abbrev=0)
APP := vis

tailwindcss:
	npx @tailwindcss/cli -i input.css -o output.css --minify

chart.js:
	curl https://cdnjs.cloudflare.com/ajax/libs/Chart.js/4.5.0/chart.umd.js > chart.js

build: chart.js
	$(MAKE) tailwindcss
	go build .

release: chart.js
	$(MAKE) tailwindcss
	GOOS=linux GOARCH=amd64 \
		go build -tags release -ldflags "-X main.version=$(VERSION)" -o $(APP)
