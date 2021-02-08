test:
	go clean -testcache
	go test ./tests/...

mediumkube:
	go build -o mediumkube main.go

mediumkubed:
	go build -o mediumkubed daemon/main.go

all: mediumkube mediumkubed

clean:
	rm -f ./mediumkube ./mediumkubed

daemon: clean mediumkubed
	sudo ./mediumkubed

