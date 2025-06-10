CGO_ENABLED=0

.PHONY: build clean

all: build

build:
	go build -trimpath -ldflags="-s -w" -o server wavycat.ru/petpet-go/cmd/app

vet:
	go vet -v ./...

test:
	go test -v ./...

clean:
	rm -f server

clean-win:
	del .\server