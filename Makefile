.PHONY: update master release update_master update_release build clean

clean:
	rm -rf vendor/
	go mod vendor

update:
	-GOFLAGS="" go get all

build:
	go build ./...
	go mod tidy

update_release:

update_master:

master: update_master clean build

release: update_release clean  build
