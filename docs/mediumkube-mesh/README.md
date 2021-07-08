# Mediumkube Mesh

## Mediumkube clustering

Some may suffer resource exaustion when launching multiple virtual machines on one single host, this is why we want to scale mediumkube horizontally to multiple host machines. Ideally, users are able to manage their virtual machines on a single node with mediumkube cli installed and access to any VM using its IP address, even if the VM could be one a different machine. Also, VMs on different hosts can communicate with each other just using their virtual IP address instead of public IP of the host machine. This is where mediumkube comes in.

## How does mediumkube mesh work?
Mediumkube mesh works pretty similarily to Kubernetes flannel network, because it utilizes flannel to establish inter-host virtual network. In overlay mode, mediumkubed is configured with CIDR, which defines the subnet where VM ips are allocated. When mediumkube mesh daemon is up, it automatically sync configurations from etcd server and modify local routing tables.

There are three cases where routing table is updated

### Lease In
Lease in means new nodes have joined the cluster since last sync. Leases in are found by looking at cidrs that are not found in local route table but found in lease definitions in etcd.

Behavior: Insert rules into routing table.

### Lease out
Opposite to Lease in condition, lease is out if the cidr is found in local routing table but nolonger exists in remote lease definitions.

Behavior: Delete rules from routing table.

### Lease outdated
Lease definition has a timestamp and a ttl defined. If timestamp + ttl is less than current timestamp, then the lease is outdated.

Behavior: Delete rules from routing table.


## How are cluster-wise domains managed

Mediumkube aggregates domains of the entire cluster by broadcasting local domains to etcd. When doing `mediumkube list`, it sends a gRPC request to local mediumkubed, and mediumkubed fetches domain list from etcd.

When creating domains, there are limitations such as unique domain name, and that comes to distributed consistency. Mediumkube implements a distributed lock using etcd, and some operations like domain deployment are locked globally. Mediumkube also has a built-in event bus that can catch up with events like domain deployed, and commerce sync event.

## Domain scheduling

When creating domains, user don't need to specify where to place it. Mediumkube automatically pick up a node in cluster and assign domain creation plan to it. 