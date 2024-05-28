resource "helm_release" "tailscale_operator" {
  name       = "tailscale-operator"
  repository = "https://pkgs.tailscale.com/helmcharts"
  chart      = "tailscale-operator"
  namespace  = "tailscale"

  set {
    name  = "oauth.clientId"
    value = var.tailscale_client_id  // Replace with your actual Tailscale client ID
  }

  set {
    name  = "oauth.clientSecret"
    value = var.tailscale_client_secret  // Replace with your actual Tailscale client secret
  }

  set {
    name  = "apiServerProxyConfig.mode"
    value = "true"
  }

  create_namespace = true
  wait             = true
}