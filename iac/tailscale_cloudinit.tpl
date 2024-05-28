#cloud-config
apt:
  sources:
    tailscale.list:
      source: deb https://pkgs.tailscale.com/stable/ubuntu focal main
      keyid: 2596A99EAAB33821893C0A79458CA832957F5868
packages:
  - tailscale
runcmd:
  - "tailscale up -authkey ${tailscale_auth_key} --advertise-routes=10.1.1.0/24,168.63.129.16/32 --accept-dns=false"
  - "echo 'net.ipv4.ip_forward = 1' | sudo tee -a /etc/sysctl.conf"
  - "echo 'net.ipv6.conf.all.forwarding = 1' | sudo tee -a /etc/sysctl.conf"
  - "sysctl -p /etc/sysctl.conf"