all: greeter.so test/test

greeter.so: greeter.go
	go build --buildmode=plugin

test/test: test/test.go
	cd test && go build

# .PHONY: proto
# proto:
# 	protoc --go_out=plugins=micro:. proto/*.proto

.PHONY: clean
clean:
	-rm -f greeter.so
	-rm -f test/test
