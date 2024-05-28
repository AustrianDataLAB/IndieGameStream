
resource azurerm_network_security_group "student" {
  name                = "student-network-security-group"
  location            = var.globals.location
  resource_group_name = data.azurerm_resource_group.rgruntime.name

  #tags = local.common_tags
}

resource "azurerm_network_security_rule" "lab_nsg" {
  name                        = "Tailscale"
  description                 = "Tailscale UDP port for direct connections. Reduces latency."
  priority                    = 1010
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Udp"
  source_port_range           = "*"
  destination_port_range      = 41641
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = data.azurerm_resource_group.rgruntime.name
  network_security_group_name = azurerm_network_security_group.student.name
}