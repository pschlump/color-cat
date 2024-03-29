
all:

build:
	go build

#install: build
#	cp color-cat ~/bin

test: build test1 test2
	@echo PASS

test1:
	./color-cat -F ./test/data2.txt >./out/test1.out
	diff ./out/test1.out ./ref/test1.ref
	@echo PASS

test2:
	./color-cat -F -n -s ./test/data2.txt >./out/test2.out
	diff ./out/test2.out ./ref/test2.ref
	@echo PASS


install: build
	rm ~/bin/color-cat
	( cd ~/bin ; ln -s ../go/src/github.com/pschlump/color-cat/color-cat . )

linux:
	GOOS=linux GOARCH=amd64 go build -o color-cat_linux

