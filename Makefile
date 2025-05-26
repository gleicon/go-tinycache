include Makefile.defs

all: deps server

deps:
	go mod tidy

server:
	go build -v -o $(NAME) cmd/main.go

clean:
	rm -f $(NAME)

test:
	go test -v
	rm bolt.db

