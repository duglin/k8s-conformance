all: build-tools doc.md tool

build-tools: bin/extract

bin/extract: src/extract.go
	@echo Building extract tool...
	go build -o $@ src/extract.go

doc.md src/funcs.go : doc.header tests/*.go src/extract.go
	@echo Building the conformance test doc...
	bin/extract tests/*.go

tool: bin/kubecon

bin/kubecon: src/kubecon.go tests/*.go src/funcs.go
	go build -o $@ src/kubecon.go src/funcs.go

clean:
	rm -f tests.md
	rm -f src/funcs.go
	rm -f bin/extract
	rm -f bin/kubecon
