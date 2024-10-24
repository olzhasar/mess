build:
	go build -o ./bin/ .

clean:
	rm -rf ./bin

.PHONY: build clean
