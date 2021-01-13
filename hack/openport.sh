#!/bin/bash

if [ -z "$1" ] || [ -z "$2" ]
  then
    echo "openport interface port"
    exit 1
fi

sudo iptables -A INPUT -p tcp --dport $2 -i $1 -j ACCEPT
sudo iptables -A OUTPUT -p tcp --dport $2 -o $1 -j ACCEPT
sudo sysctl --system