 # Go parameters
GOCMD=GO111MODULE=on GOARCH=amd64 go
GOBUILD=$(GOCMD) build -v
SOURCE_NAME=media-archiver.go
BINARY_NAME=dist/media-archiver
		
all: build-all

clean: 
	rm -rf $(BINARY_NAME)*

build: clean
	$(GOBUILD) -o $(BINARY_NAME) $(SOURCE_NAME)

build-all: clean
	GOOS=linux $(GOBUILD) -o $(BINARY_NAME)-linux $(SOURCE_NAME)
	GOOS=windows $(GOBUILD) -o $(BINARY_NAME)-windows.exe $(SOURCE_NAME)
	GOOS=darwin $(GOBUILD) -o $(BINARY_NAME)-darwin $(SOURCE_NAME)

run: build
	./$(BINARY_NAME)
