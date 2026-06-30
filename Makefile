PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin
BIN_DIR ?= bin
BIN ?= mess

.PHONY: build install clean

build:
	go build -o ./$(BIN_DIR)/$(BIN) .

install: build
	install -d "$(DESTDIR)$(BINDIR)"
	install -m 0755 "./$(BIN_DIR)/$(BIN)" "$(DESTDIR)$(BINDIR)/$(BIN)"

clean:
	rm -rf ./$(BIN_DIR)
