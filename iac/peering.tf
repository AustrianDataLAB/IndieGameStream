
## IF you have two Vnets and you want them to talk to each other, this is a bridge
## It cuts through your network segmentation, so again: be intentional about it please

# resource "azurerm_virtual_network_peering" "vnet-peering-agents" {
#   name                      = "agents_to_platform"
#   resource_group_name       = data.azurerm_resource_group.rgruntime.name
#   virtual_network_name      = azurerm_virtual_network.vnet-buildagents.name
#   remote_virtual_network_id = azurerm_virtual_network.vnet-platform.id
#   allow_virtual_network_access = true
#   allow_forwarded_traffic      = true

# }


# resource "azurerm_virtual_network_peering" "vnet-peering-platform" {
#   name                      = "platform_to_agents"
#   resource_group_name       = data.azurerm_resource_group.rgruntime.name
#   virtual_network_name      = azurerm_virtual_network.vnet-platform.name
#   remote_virtual_network_id = azurerm_virtual_network.vnet-buildagents.id
#   allow_virtual_network_access = true
#   allow_forwarded_traffic      = true

# }
