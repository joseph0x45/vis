.PHONY: tailwindcss

tailwindcss:
	tailwindcss -i input.css -o output.css

chart.js:
	curl https://cdnjs.cloudflare.com/ajax/libs/Chart.js/4.5.0/chart.umd.js > chart.js

build: chart.js
	$(MAKE) tailwindcss
	go build .
