run:
	go run main.go

build:
	go build -ldflags='-X main.VERSION=$(VERSION)' -o findr.exe

test:
	go test -v