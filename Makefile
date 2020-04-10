.PHONY: update master release setup update_master update_release build

setup:
	git config --global --add url."git@gitlab.com:".insteadOf "https://gitlab.com/"

update:
	rm -rf vendor/
	go mod vendor
	GOFLAGS="" go get -u all

build:
	go build ./...
	go mod tidy

update_release:

update_master:

master: update_master update build

release: update_release update build
