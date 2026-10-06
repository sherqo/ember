PREFIX ?= $(HOME)/.local
BINDIR = $(PREFIX)/bin

.PHONY: build install uninstall

build:
	go build -o bin/ember .

# Installs a symlink so the binary keeps its backend scripts beside it.
install: build
	mkdir -p $(BINDIR)
	ln -sf $(CURDIR)/bin/ember $(BINDIR)/ember
	@echo "ember installed → $(BINDIR)/ember"

uninstall:
	rm -f $(BINDIR)/ember
