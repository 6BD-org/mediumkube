# Setup a k8s cluster using multipass

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

## Logging

To check the log of multipass,

```bash
journalctl --unit snap.multipass*
```

If you are executing commands in the virtual machine during init, make sure to save logs to file
for analysis. In cloud-init

```bash
runcmd:
    - sh dosomething.sh >> /var/log/bootstrap/dosomething.log
```

## About CLI
For ease of development, this cli only contains 2 layers of command hierarchy and key word arguments, like following

```
$ main-command sub-command --key1 val1 --key2 val2
```

excepting `help`

```
# Get help of available sub-commands
$ main-command help 

# Get help of particular sub-command
$ main-command sub-command help
```

You can extend the layer, of course, by intercepting args and pass them through to another command handler, then you get this

```
$ main-command sub-command-1 sub-command-2 --key val
```

But I personally don't encourage either adding more layers or mixing up positional and keyword arguments.

## Roadmap
- Cli tool for cluster management
  - Cluster deployment
  - Deletion
  - Adding/Removing nodes
  - Deploy kubernetes resources

- Setup flannel network
- Better template engine