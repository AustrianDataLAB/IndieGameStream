resource "azurerm_kubernetes_cluster" "testCluster" {
  name                = "testCluster"
  location            = data.azurerm_resource_group.rgruntime.location
  resource_group_name = data.azurerm_resource_group.rgruntime.name
  dns_prefix          = "testCluster"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_B2ms"
  }

  network_profile {
    network_plugin     = "azure"
    load_balancer_sku  = "standard"
    outbound_type      = "loadBalancer"
  }

  identity {
    type = "SystemAssigned"
  }

  private_cluster_enabled = true
}

/*
output "client_certificate" {
  value = azurerm_kubernetes_cluster.testCluster.kube_config.0.client_certificate
}

output "kube_config" {
  value = azurerm_kubernetes_cluster.testCluster.kube_config_raw
}*/