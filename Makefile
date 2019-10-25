.PHONY: update master release setup update_master update_release build

setup:
	git config --global --add url."git@gitlab.com:".insteadOf "https://gitlab.com/"

update:
	rm -r vendor/
	go mod vendor
	GOFLAGS="" go get -u

build:
	go build ./...
	go mod tidy

update_release:

update_master:

master: update update_master build

release: update update_release build
