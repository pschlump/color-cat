
all:

linux:
	GOOS=linux GOARCH=amd64 go build -o color-cat_linux

deploy: linux
	scp color-cat_linux philip@${wir0dallas}:/home/philip/bin/color-cat

