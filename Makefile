
all:

build:
	go build

test: build test1 test2

test1:
	./color-cat -F ./test/data2.txt >./out/test1.out
	diff ./out/test1.out ./ref/test1.ref

test2:
	./color-cat -F -n -s ./test/data2.txt >./out/test2.out
	diff ./out/test2.out ./ref/test2.ref

