# Prerequisites
- Debian-based system
- docker
- kubectl
- kind
- helm v3.15.0+
- go v1.20.0+

You need to disable IPV6 in Docker. This can be achieved by adding adding `"ipv6": false` to the Docker daemon config.
If no file exists, this can be achieved by (**overwrites existing config file**):
```bash
echo '{"ipv6": false}' > /etc/docker/daemon.json
```
Restart the Docker daemon for the configuration to apply:
```bash
sudo systemctl restart docker
```

# How to run
## Single-node cluster
```bash
make
```

## Multi-node cluster
```bash
make multi_node
```
