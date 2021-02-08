| Parameter | Comment | Example |
|---|---|---|
|backend|virtualization backend| multipass |
|http-proxy|proxy address. Accessible via {{ .HTTPProxy }}| http://127.0.0.1:8888 |
|https-proxy|proxy address. Accessible via {{ .HTTPSProxy }}| http://127.0.0.1:8888 |
|pub-key-dir|public key for vm|./a.pub|
|priv-key-dir|private key for vm|./a|
|host-pub-key-dir|public key of host machine|./host.pub|
|image|image for vm. Only raw format is supported|./ubuntu.img|
|cloud-init|cloud init yaml file|./cloud-init.yaml|
|nodes|list of node||
|nodes.name|node name|node1|
|nodes.cpu| node cpu cores|2|
|nodes.mem|node memory size. Must be formatted as xG, where x can be an integer|2G|
|kube-init.args|list of arguments when initing kubernetes||
|kube-init.args.key|key of parameter|pod-network-cidr|
|kube-init.args.value|value of parameter|10.244.0.0/16|
|tmp_dir|working directory for mediumkube||
|vm_kube_config_dir|kube config dir in vm. You don't need to modify in most cases|/home/ubuntu/.kube/config|
