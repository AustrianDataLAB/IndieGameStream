

data "azurerm_client_config" "current" {}


#Read my target RG 

data "azurerm_resource_group" "rgruntime" {
  name     = "rg-service-not2day"

}

data "github_ip_ranges" "ranges" {}

################################################################
# IF you are generating any secrets, you need to put them somewhere
# most ideally, you put them into a keyvault of the same lifecycle-stage as the asset the key belongs to
##################################################################


 
resource "azurerm_key_vault" "kvservice" {
  name                        = "kv-service-not2day-3"
  location                    = data.azurerm_resource_group.rgruntime.location
  resource_group_name         = data.azurerm_resource_group.rgruntime.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"


  network_acls {
    ip_rules = ["0.0.0.0/0" ] #change this
    default_action= "Deny"
    bypass = "AzureServices"
  }


   access_policy {
     tenant_id = data.azurerm_client_config.current.tenant_id
     ## Students you must look up your Users Object id in Entra ID and put it here
     object_id = var.myuser



     secret_permissions = [
       "Get",
       "List",
       "Restore",
       "Delete",
       "Set",
       "Recover",
       "Backup",
     ]


   }
   # We are giving this SP access over the current vault
    access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id



    secret_permissions = [
      "Get",
      "List",
      "Restore",
      "Delete",
      "Set",
      "Recover",
      "Backup",
    ]


  }

 
tags     = local.common_tags

}

