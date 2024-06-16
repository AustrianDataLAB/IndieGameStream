
variable "subscription_id" {
  type = string
}
variable "tenant_id" {
  type = string
}



provider "azurerm" {
  subscription_id   = var.subscription_id
  tenant_id         = var.tenant_id
  storage_use_azuread = true
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "helm" {
  kubernetes {
    host                   = azurerm_kubernetes_cluster.testCluster.kube_config.0.host
    client_certificate     = base64decode(azurerm_kubernetes_cluster.testCluster.kube_config.0.client_certificate)
    client_key             = base64decode(azurerm_kubernetes_cluster.testCluster.kube_config.0.client_key)
    cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.testCluster.kube_config.0.cluster_ca_certificate)
  }
}

terraform {
  backend "azurerm" {
    use_azuread_auth = true
  }
}


locals {
  common_tags = {
    release     = "HandsOnCloudNative"
    purpose     = "class"
    classification = "sensitive"
    central     = "yes"
  }
}