#!/bin/bash
sudo systemctl stop mediumkube && make clean install && sudo systemctl start mediumkube