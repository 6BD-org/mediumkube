# Set up a k8s cluster using libvirt

It is planned that `mediunkube` no-longer depends on multipass by force, instead, multipass becomes an optional backend. Multipass is an active project, which is good, but it becomes difficult to catch up with. Therefore we move to `libvirt`, which is slower in release speed, and relatively stable and flexible.

## Prerequests

- `qemu` The hardware emulator at lowest level, which does binary translation and emulates peripheral devices
- `qemu-img` A tool used to manipulate disk images. MediumKube uses it to expand the image to desired size as user defined in yaml file
- `libvirt` libvirt is a high-level library that provides APIs for convenient manipulations of domains, networks, etc... MediumKube uses these api via rpc and some commandline tools like `virsh`, `virt-install`
- `kvm (optional)` A linux module that allows CPU to switch to guest state where privilege instructions fall back to hypervisor code. Using `kvm` along with `qemu` provides near-native performance because it avoids some unnecessary binary translations

If you have trouble installing these software, just go ahead with `multipass` backend.

## How does it work

Please refer to [this](./daemon/README.md)

## Configuration references
Please refer to [this](./docs/config.md) for configurations and [this](./docs/config-libvirt.md) for libvirt-specific configurations

## Get started
Node that mediumkube isn't a well packaged software right now, so is's not available in any kind of package manager. You will need to build the project and run binaries manually.

To compile the project

```bash
$ make clean
$ make all
```
This command will produce two executables, which are `mediumkube` and `mediumkubed`. First of all, you need to start `mediumkubed` and make sure it keeps running. Once you stop it, it will clean up configurations so you will lose ip table entries.

```bash

$ ./mediumkubed

```

Then you can deploy the machine

```bash
$ ./mediumkube deploy node1
```

This will deploy `node1` defined in your config file. The deployment process will attach you to the stdio of virtual machine, if you wanna escape, use `ctrl + ]`.

To purge the machine that is installed, use 
```bash

$ ./mediumkube purge node1

```

Please note that you cannot purge a machine that is created by another backend! So if you want to purge a multipass machine, either use `multipass` command or change backend in config file.

## Install mediumkube

Another option of using mediumkube is install it to your linux system directly

```bash
# This command will compile mediumkube and copy binarys and configurations to your system.
# The configuration files & templates are placed under /etc/mediumkube/
# The binaries are placed under /usr/local/bin
$ make clean install

```

## Use proxy

Templating engine supports proxy. So you can access `http-proxy` in your config file by using `{{ .HTTPProxy }}`. You can use any proxy, but we suggest you to deploy your proxy to listen on bridge, so that the system becomes "portable", because your nodes won't suffer from configuration changes as the network environment changes due to DHCP or switching between wifis. 

In order to set up proxy on bridge, there are two things to do. 

1. You should open port on bridge for your proxy. You can use [this script](./hack/openport.sh)
2. Just point the proxy server to the ip address of mediumkube bridge and you are good to go


## Remotely execute commands

You can execute commands on node remotely using `mediumkube`

```bash
# For example, this command lists all files under root dir
# on node1
$ mediumkube exec node1 ls /
```

## Transfer files from host to node

You can transfer files from your host machine to nodes you deployed. (Still working on another direction)

```bash
# This command sends text.txt to node1 and place it under /home/ubuntu
$ mediumkube transfer ./test.txt node1:/home/ubuntu/remote.txt
```

## Node life cycle management

In order to stop a node
```bash
$ mediumkube stop node1
```

To start a node
```bash
$ mediumkube start node1
```

To purge a node (which means stop it, then delete it along with storages attached to it)

```bash

$ mediumkube purge node1
```

## Plugins
Plugins are standalone executable units in mediumkube, they can be invoked separately with arbitrary number of arguments. The use case of plugins would be hooks integrated with main tasks such as `deploy`, `init` and `join`, etc. 

To see a list of plugins 

```bash

$ mediumkube plugin list
```

You can also see descriptions of plugins and execute them manually using `mediumkube` cli

```bash

# Execute a plugin
$ mediumkube plugin exec [plugin name] arg1 arg2 arg3

# Show descriptions of a plugin
$ mediumkube plugin desc [plugin name]
```


# [DEPRECATED] Setup a k8s cluster using multipass

This is a very simple toolkit that helps setup a K8s cluster easily (In order to learn some network knowledges about K8s)

+ Easy to use
- Unconfigurable networks
- Very simple templating
- Still need to init and join nodes manually
- No distributed deployment
- Like a minikube, but you'll have "real nodes" to access to. If you got the effort, you can config them for advanced uses

## Prepare

### Install multipass

```bash

$ sudo apt install multipass

```

In order to use multipass behind a proxy, use following command

```bash

$ sudo snap set multipass proxy.http="http://{host}:{port}"

```

### Template configurations

Most important of all, prepare three keys:
- Public key of your host machine
- Generated Private key for cluster machine
- Generated Public key for cluster machine

These are used to setup trust relations between your host and the cluster as well as cluster nodes.

When your keys are ready, modify the configuration file to point to those key files like this:

```yaml
pub-key-dir: "/home/temp/.ssh/cloud.pub"
priv-key-dir: "/home/temp/.ssh/cloud"
host-pub-key-dir: "/home/temp/.ssh/id_rsa.pub"

```

Then get your ubuntu image (`.img` file) ready, or you can simple use remote image if you are outside the bitch ass motherfucking firewall.

Also, configure the cloud-init.yaml location. It is already pointed to `./cloud-init.yaml`, which is the default output of template renderer. If you change this, make sure it exists.

```yaml
image: "file:///home/temp/u_20.04.img"
cloud-init: "./cloud-init.yaml"
```

Finally, if you need proxy, do configurations like this

```yaml
http-proxy: "http://localhost:1091"
https-proxy: "http://localhost:1091"
```
and use `{{ .HTTPSProxy }}` to configure your software.

However, if no proxy is required, remember to remote related template tokens from .tmpl file. This may take some effort :smirk:

Now you are ready to go, build the project and setup your cluster

### Test & Build

[Golang officially](https://golang.org/pkg/testing/) suggests to put test files together with bussiness logic, but we have too many mock-data files, so note that ALL unit tests are located in `./tests`. To run test, 

```bash

$ go test ./tests/...

```

In order to build the project, 

```bash

$ ./hack/build.sh

```

This will generate an executable `main` in project root.


## Templating Guide

In order to simplify the configuration, we support configuration
and template rendering

there are pre-build options which are proxies. 

```yaml
http-proxy: "172.16.184.20:1091"
https-proxy: "172.16.184.20:1091"
```

In order to use configuration instead of writing proxies everywhere, use 

```
{{ .HttpProxy }} and {{ .HttpsProxy }}
```

Also be careful when processing data sensitive fields like private key. Using go template might introduce one `newline` to template file, so remember to trim. 

```yaml
privKey: |
    {{ .PrivKey | nindent 6 }}
```
For example this is translated to 

```yaml
privKey: |
# Note there's a newline below, thus the key is incorrect

    -----BEGIN RSA PRIVATE KEY-----
    asdasdasdasdasdasd....
```

Instead you do this

```yaml
privKey: |
    {{- PrivKey | nindent 6 }}
```

in your `yaml.tmpl` file, and render it using that simple go program

```bash

$ ./main render

```

To get help, of available commands

```bash
# List available commands
$ ./main help

# Get help of sub commands
$ ./main render help

```

## Multipass compatibility
This cli is fully compatible with multipass. You can replace `multipass` with `./main` or any executable name that you build. Just for consistent-looking :smirk:

```bash

# These commands are identical

$ multipass list

$ ./main list 

```

Checkout the [multipass documentation](https://multipass.run/docs/working-with-instances)

## Launch instance

```bash

# -c 2 uses 2 cpus
# -m 2G 2G memory
# -d 20G 20G disk
# -n node01 node named node01
# file path to .img file

$ multipass launch -v -n node01 --cloud-init cloud-init.yaml -c 2 -m 2G -d 20G file:///home/temp/u_20.04.img
```

A better way of launching instance is via cli

```bash

$ ./main deploy --config ./cloud-init.yaml
```

## purge instance

```bash

$ ./hack/purge.sh {instance_name}

# To purge multiple nodes at the same time
$ ./hack/purge.sh node1 node2 node3
```

## Start K8s cluster


```
# We normally don't have enough resource
# for launching cluster so add this 
# flag
# --ignore-preflight-errors=all


# To start a master node, do this on node01
$ kubeadm init --ignore-preflight-errors=all
```

A better way of starting k8s cluster is using out cli after 
configuring kube-init section properly.

```
$ ./main init --config ./config.yaml
```


## Install resource to kubernetes using MediumKube

Following types are currently supported You are free to add more if you need them

```golang
	resourceType["PodSecurityPolicy"] = &v1beta1.PodSecurityPolicy{}
	resourceType["ClusterRole"] = &v1.ClusterRole{}
	resourceType["ClusterRoleBinding"] = &v1.ClusterRoleBinding{}
	resourceType["ServiceAccount"] = &coreV1.ServiceAccount{}
	resourceType["ConfigMap"] = &coreV1.ConfigMap{}
	resourceType["DaemonSet"] = &appsV1.DaemonSet{}
	resourceType["StatefulSet"] = &appsV1.StatefulSet{}
```

You can edit your yaml outside the cluster using your favorite text editor, and submit them using the command 

```bash
$ ./mediumkube apply my.yaml
```


## Roadmap
- Cli tool for cluster management
  - Cluster deployment
  - Deletion
  - Adding/Removing nodes
  - Deploy kubernetes resources

- Setup flannel network
- Better template engine