# Setup a k8s cluster using multipass

## Prepare

### Install multipass

```bash

$ sudo apt install multipass

```

### Launch instance

```bash
$ multipass launch -v -n node01 --cloud-init init.yaml file://{path_to_image}
```

### purge instance

```bash

$ ./purge.sh {instance_name}
```
