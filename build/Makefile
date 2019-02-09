GOBIN = $(shell pwd)
GO ?= latest

GO_BUILD_EX_ARGS ?=

all: openx emulator teller test

.PHONY: openx emulator teller test tests clean

goget:
	go get github.com/stellar/go
	# stellar/go problems, so need to have dep
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	cd ${GOBIN}/gopath/src/github.com/stellar/go
	build/env.sh go get -v ./...

openx: goget
	build/env.sh go build ${GO_BUILD_EX_ARGS}
	@echo "Run \"$(GOBIN)/openx\" to launch openx."

emulator: goget
	build/env.sh go build ${GO_BUILD_EX_ARGS} ./emulator/emulator
	@echo "Run \"$(GOBIN)/emulator\" to launch emulator."

teller: goget
	build/env.sh go build ${GO_BUILD_EX_ARGS} ./teller/teller
	@echo "Run \"$(GOBIN)/teller\" to launch teller."

clean:
	build/env.sh clean
	go clean .
	go clean ./emulator/emulator
	rm -rf build/_workspace/
	rm -f openx emulator/emulator teller/teller

test tests: openx
	build/env.sh go test -v ./...
