MAKEPATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
NAME:=statusctl

help:
	@echo View the Makefile targets for usage.

build: clean
	cd $(MAKEPATH); go build -o $(NAME) .

run: build
	cd $(MAKEPATH); ./$(NAME)

clean:
	cd $(MAKEPATH); go fmt ./...
	cd $(MAKEPATH); go vet ./...

tidy:
	cd $(MAKEPATH); go mod tidy
