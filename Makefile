GOOS=linux
GOARCH=amd64

.PHONY: build
build: clean
	go build -o ./dist/step-functions-agent main.go

.PHONY: clean
clean:
	rm -rf dist

