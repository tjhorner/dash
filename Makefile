.PHONY: dist dist-win dist-macos dist-linux ensure-dist-dir build build-frontend install uninstall

GOBUILD=packr2 build -ldflags="-s -w"
INSTALLPATH=/usr/local/bin

ensure-dist-dir:
	@- mkdir -p dist

build-frontend:
	cd frontend && npm run build

dist-win: ensure-dist-dir
	# Build for Windows x64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o dist/dash-windows-amd64.exe *.go
	packr2 clean

dist-macos: ensure-dist-dir
	# Build for macOS x64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o dist/dash-darwin-amd64 *.go
	packr2 clean

dist-linux: ensure-dist-dir
	# Build for Linux x64
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o dist/dash-linux-amd64 *.go
	packr2 clean

dist: build-frontend dist-win dist-macos dist-linux clean

clean:
	packr2 clean

build: build-frontend
	@- mkdir -p bin
	$(GOBUILD) -o bin/dash *.go
	make clean
	@- chmod +x bin/dash

install: build
	mv bin/dash $(INSTALLPATH)/dash
	@- rm -rf bin
	@echo "dash was installed to $(INSTALLPATH)/dash. Run make uninstall to get rid of it, or just remove the binary yourself."

uninstall:
	rm $(INSTALLPATH)/dash

run:
	@- go run *.go