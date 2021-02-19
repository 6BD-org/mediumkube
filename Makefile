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

install: mediumkube mediumkubed
	sudo cp mediumkube /usr/local/bin/mediumkube
	sudo cp mediumkubed /usr/local/bin/mediumkubed
	sudo mkdir -p /etc/mediumkube
	sudo cp config.yaml /etc/mediumkube/config.yaml

