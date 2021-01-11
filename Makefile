mediumkube:
	go build -o mediumkube main.go

mediumkubed:
	go build -o mediumkubed daemon/main.go

all: mediumkube mediumkubed

clean:
	rm ./mediumkube ./mediumkubed