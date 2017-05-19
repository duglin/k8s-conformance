all: build-tools doc.md tool

build-tools: bin/extract

bin/extract: src/extract.go
	echo Building extract tool...
	go build -o $@ src/extract.go

doc.md src/funcs.go : doc.header tests/*.go src/extract.go
	echo Building the conformance doc...
	cat doc.header > doc.md
	echo >> doc.md
	bin/extract tests/*.go
	cat doc.body >> doc.md
	rm doc.body

tool: bin/k8scon

bin/k8scon: src/k8scon.go tests/*.go src/funcs.go
	go build -o $@ src/k8scon.go src/funcs.go

clean:
	rm -f doc.mg
	rm -f bin/extract
	rm -f bin/k8scon
	rm -f src/funcs.go
	rm -f doc.md
