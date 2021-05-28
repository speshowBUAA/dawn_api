# Go parameters 
GOCMD=go 
GOBUILD=$(GOCMD) build 
GOCLEAN=$(GOCMD) clean 
GOTEST=$(GOCMD) test 
GOGET=$(GOCMD) get 

GOOS=linux
GOARCH=amd64
#GOARCH=arm
BINARY_NAME=dawn_api

BINARY_UNIX=$(BINARY_NAME)_unix 
# export GOPATH=$(GOPATH):.

# CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' .
all: test build 
build: 
	env GOOS=$(GOOS) CGO_ENABLED=0 GOARCH=$(GOARCH) $(GOBUILD) -o $(BINARY_NAME) -v -a -ldflags '-extldflags "-static"' 
test:
#	$(GOTEST) -v ./...
	echo $(GOPATH)
clean: 
	$(GOCLEAN) 
	rm -f $(BINARY_NAME) 
	rm -f $(BINARY_UNIX) 
run: 
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)
