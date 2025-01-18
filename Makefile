start: build
	@./telepath daemon start

stop: build
	@./telepath daemon stop

status: build
	@./telepath daemon status

build:
	@go build -o telepath .
