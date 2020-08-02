BINARY_NAME=server

all: prod fmt build

build:
	go build -o bin/$(BINARY_NAME) -v main.go

run:
	./bin/$(BINARY_NAME) -f etc/config.yml

clean:
	rm -f bin/$(BINARY_NAME)

fmt:
	go fmt ./...

prod:
	npm run prod