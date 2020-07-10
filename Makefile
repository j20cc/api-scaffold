BINARY_NAME=server

all: fmt prod build

build:
	go build -o bin/$(BINARY_NAME) -v app/main.go

run:
	./bin/$(BINARY_NAME) -f config.yml

clean:
	rm -f bin/$(BINARY_NAME)

fmt:
	go fmt ./...

prod:
	npm run prod