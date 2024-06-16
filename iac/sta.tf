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
    default_action = "Deny"
    ip_rules = ["0.0.0.0/0"]
    bypass = ["AzureServices"]
  }
}


resource "azapi_resource" "gamesContainer" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-01-01"
  name      = "games"
  parent_id = "${azurerm_storage_account.staindiegamestream.id}/blobServices/default"
  body = jsonencode({
    properties = {
      publicAccess = "None"
    }
  })
  depends_on = [
    azurerm_storage_account.staindiegamestream
  ]
}

resource "azurerm_role_assignment" "role_assignment" {
  scope                  = azurerm_storage_account.staindiegamestream.id
  role_definition_name   = "Storage Blob Data Contributor"
  principal_id           = data.azurerm_client_config.current.object_id
}