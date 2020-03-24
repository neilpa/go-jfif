all: jfifstat xmpdump

jfifstat: bin/jfifstat
xmpdump: bin/xmpdump

bin/jfifstat: *.go cmd/jfifstat/*.go
	go build -o $@ ./cmd/jfifstat
bin/xmpdump: *.go cmd/xmpdump/*.go
	go build -o $@ ./cmd/xmpdump

test:
	go test ./...

fmt:
	go fmt ./...
tidy:
	go mod tidy

clean:

.PHONY: all jfifstat xmpdump test fmt tidy clean
