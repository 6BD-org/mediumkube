SYSTEMD_DIR=/lib/systemd/system

test:
	go clean -testcache
	go test ./tests/...

mediumkube:
	go build -o build/mediumkube main.go

mediumkubed:
	go build -o build/mediumkubed pkg/daemon/main.go

all: mediumkube mediumkubed

clean:
	rm -f build/*

install: mediumkube mediumkubed
	sudo mkdir -p /etc/mediumkube /var/run/mediumkube
	sudo mkdir -p /etc/mediumkube/flannel /var/run/mediumkube/flannel

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

stop:
	sudo systemctl stop mediumkube

start:
	sudo systemctl start mediumkube	

