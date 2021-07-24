SYSTEMD_DIR=/lib/systemd/system
GRPC_ROOT=pkg/daemon/mgrpc


test:
	go clean -testcache
	go test ./tests/...

generate:
	protoc --go_out=. \
	--go_opt=paths=source_relative \
	 --go-grpc_out=. \
	 --go-grpc_opt=paths=source_relative \
	 daemon/mgrpc/domain.proto

mediumkube:
	go build -o build/mediumkube commands/main.go

mediumkubed:
	go build -o build/mediumkubed daemon/main.go

all: mediumkube mediumkubed

clean:
	rm -f build/*

install: mediumkube mediumkubed
	sudo mkdir -p /etc/mediumkube /var/run/mediumkube
	sudo mkdir -p /etc/mediumkube/flannel /var/run/mediumkube/flannel

	sudo cp context/* /usr/local/bin

	# Copy binary and default configuration files
	sudo cp build/mediumkube /usr/local/bin/mediumkube
	sudo cp build/mediumkubed /usr/local/bin/mediumkubed
	sudo cp config.yaml /etc/mediumkube/config.yaml
	
	# Register systemd service
	sudo cp mediumkube.service.start.sh /usr/local/sbin && sudo chmod +x /usr/local/sbin/mediumkube.service.start.sh
	sudo cp mediumkube.service.stop.sh /usr/local/sbin && sudo chmod +x /usr/local/sbin/mediumkube.service.stop.sh
	sudo cp mediumkube.service $(SYSTEMD_DIR)/mediumkube.service
	sudo systemd-analyze verify $(SYSTEMD_DIR)/mediumkube.service

	# Reload and enable service
	sudo systemctl daemon-reload
	sudo systemctl enable mediumkube

.PHONY: stop
stop:
	sudo systemctl stop mediumkube

.PHONY: start
start:
	sudo systemctl start mediumkube	

