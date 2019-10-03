DEP := $(shell command -v dep 2> /dev/null)
PACKAGES=$(shell go list ./... | grep -v '/vendor/')

export GO111MODULE = on

build:
	go build  -o ./build
	go build -mod vendor -o build