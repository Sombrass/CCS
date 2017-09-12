SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

test:
	go test $(TARGET)

fmt:
	@gofmt -l -w $(SRC)


