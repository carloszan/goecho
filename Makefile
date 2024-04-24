run: build
	./bin/goecho --listenAddr :5001

build:
	go build -o bin/goecho .
