These are libvirt-specific configurations

| Parameter | Comment | Example |
|---|---|---|
|bridge|bridge configuration for mediumkube||
|bridge.name|name of the bridge. No need to change in most cases|mediumkube-br0|
|bridge.inet|ipv4 address for bridge. Must include netmask. Change this if there are conflicts with existing interfaces on your host|10.114.114.1/24|
|bridge.host|host NIC. Not used for now|enp1|
