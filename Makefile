BINARY=cargo

BUILD=$$(vtag)

REVISION=`git rev-list -n1 HEAD`

default: build

build: format
	go build -o ${BINARY} -v *.go

format fmt:
	gofmt -l -w .

clean:
	go mod tidy
	go clean
	rm $(BINARY)

get-tag:
	echo ${BUILD}

.PHONY: build format fmt clean get-tag
