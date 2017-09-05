all: frame

# proto:
# 	for f in exts/proto/**/*.proto; do \
# 		protoc --go_out=plugins=micro:. $$f; \
# 		echo compiled: $$f; \
# 	done \

# SRC = $(wildcard exts/*.go)
# OBJ = $(patsubst %.go,plugins/%.so,$(notdir ${SRC}))
# plugins/%.so:exts/%.go
# 	mkdir -p plugins/
# 	go build --buildmode=plugin -o $@ $<

plugin: external/
	make -C external/

frame: */*.go *.go ./internal/internal.go
	go build

./internal/internal.go: ./internal
	cd internal && python3 mkinternal.py | tee internal.go

.PHONY: install
install:
	mkdir -p plugins
	cp external/*/*.so plugins/

# test:
# 	./microframe srv -c ./conf/config.yml

.PHONY: clean
clean:
	-rm -f plugins/*.so
	-rm -f microframe
