run:
	go run main.go

build:
	go build -ldflags='-X main.VERSION=$(VERSION)' -o findr.exe && GOOS=linux GOARCH=amd64 go build -ldflags='-X main.VERSION=$(VERSION)' -o findr

test:
	go test -v