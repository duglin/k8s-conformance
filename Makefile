all: build-tools doc.md tool

BUILD_OPTS=src/kubecon.go src/funcs.go

build-tools: bin/extract

bin/extract: src/extract.go
	@echo Building extract tool...
	go build -o $@ src/extract.go

doc.md src/funcs.go : tests/*.go src/extract.go
	@echo Building the conformance test doc...
	bin/extract tests/*.go

tool: bin/kubecon

bin/kubecon: src/kubecon.go tests/*.go src/funcs.go
	go build -o $@ ${BUILD_OPTS}

cross:
	BUILD_OPTS="${BUILD_OPTS}" utils/cross.sh

clean:
	rm -f tests.md
	rm -f src/funcs.go
	rm -f bin/extract
	rm -f bin/kubecon
