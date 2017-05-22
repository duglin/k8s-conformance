all: build-tools tests.md tool verify

BUILD_OPTS=-tags netgo -installsuffix netgo src/kubecon.go src/funcs.go
DO_UPDATE=
ifeq ("$(wildcard bin/verify*)","")
DO_UPDATE=update
endif

build-tools: bin/extract

bin/extract: src/extract.go
	@echo Building extract tool...
	go build -o $@ src/extract.go

tests.md src/funcs.go : tests/*.go bin/extract
	@echo Building the conformance test doc...
	bin/extract tests/*.go

tool: bin/kubecon

bin/kubecon: src/kubecon.go src/funcs.go
	go build -o $@ ${BUILD_OPTS}

cross:
	BUILD_OPTS="${BUILD_OPTS}" utils/cross.sh

verify: $(DO_UPDATE)
	bin/verify-links.sh *.md

update:
	curl -s https://raw.githubusercontent.com/duglin/vlinker/master/bin/verify-links.sh > bin/verify-links.sh
	chmod +x bin/verify-links.sh

clean:
	rm -f tests.md
	rm -f src/funcs.go
	rm -f bin/kubecon* bin/extract

purge: clean
	rm  -rf bin
