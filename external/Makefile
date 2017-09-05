plugins: family # single

# SINGLE_SRC = $(wildcard *.go)
# SINGLE_SO = $(patsubst %.go,../plugins/%.so,$(notdir ${SINGLE_SRC}))
# single:${SINGLE_SO}
# ../plugins/%.so:%.go
# 	go build --buildmode=plugin -o $@ $<

# EXCLUDE_DIR := proto
EXTDIR := $(shell find . -maxdepth 1 -type d)
EXTDIR := $(basename $(patsubst ./%,%,$(EXTDIR)))
# EXTDIR := $(filter-out $(EXCLUDE_DIR),$(EXTDIR))
FAMILY_SO := $(patsubst %,%/_tmp.so,$(EXTDIR))
# FAMILY_SO := $(foreach n,$(EXTDIR),$(n)/$(n).so)

family: $(FAMILY_SO)

%/_tmp.so: %
	echo $@ $<
	if [ -f $</Makefile ]; then \
		make -C $<; \
	else \
		go build --buildmode=plugin -o $@ ./$<; \
		mv $@ $</$<.so; \
	fi

