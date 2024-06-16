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
}

resource "azurerm_storage_container" "gamesContainer" {
  name                      = "games"
  storage_account_name      = azurerm_storage_account.staindiegamestream.name
  container_access_type     = "blob"
}

resource "azurerm_storage_account_network_rules" "netrules" {
  storage_account_id = azurerm_storage_account.staindiegamestream.id

  default_action = "Deny"
  bypass = ["AzureServices"]

  depends_on = [
    azurerm_storage_container.gamesContainer,
  ]
}