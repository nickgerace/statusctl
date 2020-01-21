MAKEPATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
NAME:=statusctl

help:
	@echo View the Makefile targets for usage.

# Development targets.

build: clean
	cd $(MAKEPATH); go build -o $(NAME) .

run: build
	cd $(MAKEPATH); ./$(NAME)

clean:
	cd $(MAKEPATH); go fmt ./...
	cd $(MAKEPATH); go vet ./...

tidy:
	cd $(MAKEPATH); go mod tidy

# Installation targets.

install:
	sudo mv $(MAKEPATH)/$(NAME) /usr/local/bin/

uninstall:
	if [ -x /usr/local/bin/$(NAME) ]; then rm -ri $(HOME)/.config/$(NAME); fi
	if [ -x /usr/local/bin/$(NAME) ]; then sudo rm /usr/local/bin/$(NAME); fi
