VERSION := $(or $(AppVersion), "v0.0.0")
COMMIT := $(or $(shell git rev-parse --short HEAD), "unknown")
BUILDDATE := $(shell date +%Y-%m-%d)

LDFLAGS := -X 'main.AppVersion=$(VERSION)' -X 'main.CommitHash=$(COMMIT)' -X 'main.BuildDate=$(BUILDDATE)'

all: build

setup:
	go mod tidy

test:
	go test -v ./...  -race -coverprofile=coverage.out -covermode=atomic

coverage: test
	go tool cover -func=coverage.out

coverage-html: test
	mkdir -p coverage
	go tool cover -html=coverage.out -o coverage/index.html

coverage-serve: coverage-html
	python3 -m http.server 8080 -d coverage

install: build
	cp telepath /usr/local/bin/telepath
	cp man/telepath.1 /usr/local/share/man/man1/telepath.1

uninstall:
	rm /usr/local/bin/telepath
	rm /usr/local/share/man/man1/telepath.1

build:
	go build -gcflags="all=-N -l" -ldflags="$(LDFLAGS)" -o telepath

dist:
	cp man/telepath.1 man/telepath.old
	sed -e "s|BUILDDATE|$(BUILDDATE)|g" -e "s|VERSION|$(VERSION)|g" man/telepath.old > man/telepath.1

	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/telepath-linux-amd64
	cp build/telepath-linux-amd64 build/telepath
	tar -zcvf build/telepath-linux-amd64.tar.gz build/telepath man/telepath.1

	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o build/telepath-linux-arm64
	cp build/telepath-linux-arm64 build/telepath
	tar -zcvf build/telepath-linux-arm64.tar.gz build/telepath man/telepath.1

	GOOS=linux GOARCH=arm go build -ldflags="$(LDFLAGS)" -o build/telepath-linux-arm
	cp build/telepath-linux-arm build/telepath
	tar -zcvf build/telepath-linux-arm.tar.gz build/telepath man/telepath.1

	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/telepath-darwin-amd64
	cp build/telepath-darwin-amd64 build/telepath
	tar -zcvf build/telepath-darwin-amd64.tar.gz build/telepath man/telepath.1

	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o build/telepath-darwin-arm64
	cp build/telepath-darwin-arm64 build/telepath
	tar -zcvf build/telepath-darwin-arm64.tar.gz build/telepath man/telepath.1
	rm build/telepath

	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/telepath-windows-amd64.exe
	GOOS=windows GOARCH=386 go build -ldflags="$(LDFLAGS)" -o build/telepath-windows-i386.exe

	# Generating checksum
	cd build && sha256sum * >> checksum-sha256sum.txt
	cd build && md5sum * >> checksum-md5sum.txt

	# Cleaning
	mv man/telepath.old man/telepath.1

clean:
	rm -rf telepath*
	rm -rf build

# For headless debugging
debug-srv-headless: build
	dlv exec telepath --headless --listen=:2345 --api-version=2 -- daemon start --daemon-child

# Will connect remote debugger
debug-connect:
	dlv connect :2345

# Will debug daemon locally
debug: build
	dlv exec telepath -- daemon start --daemon-child
# For client debugging you need to start similar command line this
# dlv exec telepath -- daemon status
# dlv exec telepath -- host list
