include Makefile.defs

all: deps server

deps:
	go mod tidy

server:
	go build -v -o $(NAME) cmd/go-tinycache/main.go

clean:
	rm -f $(NAME)

test:
	go test -v
	rm bolt.db

