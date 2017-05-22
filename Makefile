all: build-tools tests.md tool verify

BUILD_OPTS=-tags netgo -installsuffix netgo src/kubecon.go src/funcs.go
DO_UPDATE=
ifeq ("$(wildcard utils/verify*)","")
DO_UPDATE=update
endif

build-tools: bin/extract

bin/extract: src/extract.go
	@echo -e \\nBuilding extract tool...
	go build -o $@ src/extract.go

tests.md src/funcs.go : tests/*.go bin/extract
	@echo -e \\nBuilding the conformance test doc...
	bin/extract tests/*.go

tool: bin/kubecon

bin/kubecon: src/kubecon.go src/funcs.go
	go build -o $@ ${BUILD_OPTS}

cross:
	BUILD_OPTS="${BUILD_OPTS}" utils/cross.sh

verify: $(DO_UPDATE)
	@echo -e \\nRunning the href checking tool...
	utils/verify-links.sh *.md

update:
	@echo -e \\nDownloading the href checking tool...
	curl -s https://raw.githubusercontent.com/duglin/vlinker/master/bin/verify-links.sh > utils/verify-links.sh
	chmod +x utils/verify-links.sh

clean:
	rm -f tests.md
	rm -f src/funcs.go
	rm -rf bin

purge: clean
	rm  -f utils/verify-links.sh
