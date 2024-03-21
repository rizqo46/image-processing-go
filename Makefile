run:
	go run main.go

build:
	go build

build-and-run:
	go build 
	./image-processing-go

test:
	go test ./... -cover

test-view-html:
	go test ./... -coverprofile=c.out
	go tool cover -html="c.out"