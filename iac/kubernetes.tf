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

  storage_profile {
    blob_driver_enabled = true
  }

  identity {
    type = "SystemAssigned"
  }

  private_cluster_enabled = true
}

resource "azurerm_storage_account" "staindiegamestream" {
  name                            = "staindiegamestream"
  resource_group_name             = data.azurerm_resource_group.rgruntime.name
  location                        = data.azurerm_resource_group.rgruntime.location
  
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  account_kind                    = "StorageV2"

  public_network_access_enabled   = false

  nfsv3_enabled                   = true
  is_hns_enabled                  = true

  network_rules {
    default_action                = "Deny"
    ip_rules                      = ["0.0.0.0/0"]
    bypass                        = ["AzureServices"] 
  }
}

/*
output "client_certificate" {
  value = azurerm_kubernetes_cluster.testCluster.kube_config.0.client_certificate
}

output "kube_config" {
  value = azurerm_kubernetes_cluster.testCluster.kube_config_raw
}*/