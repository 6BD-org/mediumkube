# MediumKubed

MediumKubed is daemon for mediumkube. This is introduced to support `libvirt` backend. 

## What does mediumkubed does?
- Maintain a virtual bridge network
- Monitor and configure virtual bridge
- Other tasks 


## How does it work?

In order to make out system working, it is preferred that we can have a virtual network controlled by mediumkube. Therefore we introduce a daemon called `mediumkubed`, that automatically configures virtual network and iptable entries for us. The logic of `mediumkubed` is like this 
![](./mediumkubed-design.png)


## Network Topology

The network topology looks like this 
![](./network-design.png)