# Mediumkube Mesh

## Mediumkube clustering

Some may suffer resource exaustion when launching multiple virtual machines on one single host, this is why we want to scale mediumkube horizontally to multiple host machines. Ideally, users are able to manage their virtual machines on a single node with mediumkube cli installed and access to any VM using its IP address, even if the VM could be one a different machine. Also, VMs on different hosts can communicate with each other just using their virtual IP address instead of public IP of the host machine. This is where mediumkube comes in.