build:
	go build -o ./bin/ .

install: build
	cp ./bin/* ~/.local/bin/

clean:
	rm -rf ./bin

.PHONY: build clean
