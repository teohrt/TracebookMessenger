GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=TracebookMessenger

build:
	$(GOBUILD)

run-first: build
	./$(BINARY_NAME) -name=First -port=1111

run-second: build
	./$(BINARY_NAME) -name=Second -port=2222 -peerAddress=192.168.0.40:1111

run-third: build
	./$(BINARY_NAME) -name=Third -port=3333 -peerAddress=192.168.0.40:2222

clean:
	rm $(BINARY_NAME)