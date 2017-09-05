all: greeter.so test/client

greeter.so: greeter.go
	go build --buildmode=plugin

test/client: test/client.go
	cd test && go build

# .PHONY: proto
# proto:
# 	protoc --go_out=plugins=micro:. proto/*.proto

.PHONY: clean
clean:
	-rm -f greeter.so
	-rm -f test/client
