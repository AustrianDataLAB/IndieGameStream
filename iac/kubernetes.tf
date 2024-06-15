resource "azurerm_kubernetes_cluster" "testCluster" {
  name                = var.cluster_name
  location            = data.azurerm_resource_group.rgruntime.location
  resource_group_name = data.azurerm_resource_group.rgruntime.name
  dns_prefix          = var.cluster_name

  default_node_pool {
    name       = "default"
    node_count = 2
    vm_size    = "Standard_B2ms"
    upgrade_settings {
      drain_timeout_in_minutes = 5
      max_surge = "50%"
      node_soak_duration_in_minutes = 0
    }
    max_pods = 110
    temporary_name_for_rotation = "upgrade"
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