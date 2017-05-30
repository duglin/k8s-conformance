# Targets of interest:
#   all - build docs and 'kubecon' tool for current platform/OS, verify hrefs
#   cross - builds docs, 'kubecon' tool for all platforms/OS

EXE=bin/kubecon
BUILD_OPTS=-tags netgo -installsuffix netgo src/kubecon.go src/funcs.go

all: $(EXE) verify

bin/extract: src/extract.go utils/test.go
	@echo -e \\nBuilding extract tool...
	go build -o $@ src/extract.go

tests.md src/funcs.go : tests/*.go bin/extract
	@echo -e \\nBuilding the conformance test doc...
	@#Tests are run in the order in which the list are listed on the
	@#next comment and in the order specified within the golang file
	bin/extract tests/rs.go tests/pod.go

$(EXE): src/kubecon.go src/funcs.go utils/*.go
	go build -o $@ ${BUILD_OPTS}

cross: .cross
.cross: $(EXE)
	@echo -e \\nBuilding \'kubecon\' for all platforms...
	@BUILD_OPTS="${BUILD_OPTS}" BINARY="${EXE}" utils/cross.sh
	@echo > .cross

verify: .verify
.verify: tests.md utils/verify-links.sh
	@echo -e \\nRunning the href checking tool...
	utils/verify-links.sh *.md
	@echo > .verify

utils/verify-links.sh:
	@echo -e \\nDownloading the href checking tool...
	curl -s https://raw.githubusercontent.com/duglin/vlinker/master/bin/verify-links.sh > $@
	chmod +x $@

clean:
	rm -rf tests.md src/funcs.go .cross .verify bin

purge: clean
	rm -f utils/verify-links.sh
