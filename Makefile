build:
	go build -o ./bin/ .

install: build
	cp ./bin/* ~/.local/bin/

clean:
	rm -rf ./bin

test:
	go test ./...

.PHONY: build clean test
