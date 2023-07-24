.PHONY: update master release update_master update_release build clean

clean:
	go mod tidy
	go mod vendor -e

update:
	-GOFLAGS="" go get all

build:
	go build ./...

update_release:

update_master:

master: update_master clean build

release: update_release clean build
