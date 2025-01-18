start: build
	./telepath start --daemon

stop: build
	./telepath kill --daemon

build:
	go build -o telepath .
