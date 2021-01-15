# Virtualization on Fedora

```sh
# List related packages in virtualization package group
$ dnf groupinfo virtualization
```

```sh
# Install packages
$ sudo dnf install @virtualization
```


```sh
# Start libvirtd
$ sudo systemctl start libvirtd
```

```sh
# Then you can use packages
$ virsh list
```