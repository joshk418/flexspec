.PHONY: build-ui build test

build-ui:
	cd ui && npm ci && npm run build

build: build-ui
	go build -o flexspec .

test:
	go test -race ./...
