default: build

build:
	go build -o ./bin/ .

install: build
	install -m 0755 ./bin/* ~/.local/bin/

clean:
	rm -rf ./bin

test:
	go test ./...

test_watch:
	rg --files -g '*.go' -g justfile | entr -c just test

fmt:
	gofmt -w .

vet:
	go vet ./...

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

check: fmt vet test

run *args:
	go run . {{args}}
