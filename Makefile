.PHONY: build serve engine client container tidy pushgit fmt generate

GO = /usr/local/go/bin/go

build:
	cd bitkv-ctl && $(GO) build -o $$GOPATH/bin/bitkv-ctl
	cd bitkv-server && $(GO) build -o $$GOPATH/bin/bitkv-server

fmt:
	gofmt -w .

test:
	$(GO) test ./...
