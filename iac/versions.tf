terraform {
  required_providers {

    azuread = {
      source  = "hashicorp/azuread"
      version = ">= 2.7.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 2.59.0"
    }
    tls = {
      source = "hashicorp/tls"
      version = "4.0.4"
    }
    azapi = {
      source = "Azure/azapi"
      version = "1.13.1"
    }
  }
  required_version = ">= 0.13"
}