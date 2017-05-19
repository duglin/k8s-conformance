all: tools doc.md

tools: bin/extract

bin/extract: src/extract.go
	echo Building extract tool...
	go build -o $@ src/extract.go

doc.md: doc.header tests/*.go
	echo Building the conformance doc...
	cat doc.header > doc.md
	echo >> doc.md
	bin/extract tests/*.go >> doc.md
