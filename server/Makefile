GOFMT=gofmt
GC=go build
VERSION := $(shell git describe --always --tags --long)
BUILD_NODE_PAR = -ldflags "-X git.ont.io/Onchain/CandyDrop/config.Version=$(VERSION)" #-race

candybox-cross: candybox-windows candybox-linux candybox-darwin

candybox-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o candybox-windows-amd64.exe main.go

candybox-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o candybox-linux-amd64 main.go

candybox-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o candybox-darwin-amd64 main.go