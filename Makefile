CC=go build
MAC=GOOS=darwin GOARCH=amd64
LIN64=GOOS=linux GOARCH=amd64
LIN32=GOOS=linux GOARCH=386

all: run

run: go run .

build-all: mac lin64 lin32

mac:
	$(MAC) $(CC) -o ./bin/proxy-runner-mac

lin64:
	$(LIN64) $(CC) -o ./bin/proxy-runner-lin64

lin32:
	$(LIN32) $(CC) -o ./bin/proxy-runner-lin32

clean:
	rm -rf ./bin/proxy-runner*

rebuild: clean build
