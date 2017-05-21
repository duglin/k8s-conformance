all: build-tools tests.md tool

BUILD_OPTS=-tags netgo -installsuffix netgo src/kubecon.go src/funcs.go

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

clean:
	rm -f tests.md
	rm -f src/funcs.go
	rm -rf bin
