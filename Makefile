BINARY_NAME=server

all: prod fmt build

build:
	go build -o bin/$(BINARY_NAME) -v cmd/main.go

run:
	./bin/$(BINARY_NAME)

clean:
	rm -f bin/$(BINARY_NAME)

fmt:
	go fmt ./...

prod:
	npm run prod