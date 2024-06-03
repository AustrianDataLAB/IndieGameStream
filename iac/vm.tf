
/*

resource "azurerm_network_interface" "student" {
  name                = "nic-student"
  location            = data.azurerm_resource_group.rgruntime.location
  resource_group_name = data.azurerm_resource_group.rgruntime.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     =  azurerm_subnet.snet-student-vm.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "tls_private_key" "example_ssh" {
    algorithm = "RSA"
    rsa_bits = 4096
}

## Commented out to save money , VMs cost real money! Dont leave them on if you dont need them



resource "azurerm_linux_virtual_machine" "vm" {
  name                = "vm"
  resource_group_name = data.azurerm_resource_group.rgruntime.name
  location            = data.azurerm_resource_group.rgruntime.location
  size                = "Standard_F1"
  priority            = "Spot"
  eviction_policy     = "Deallocate"
  disable_password_authentication = "true"
  admin_username      = "adminusercrcr"
  network_interface_ids = [
    azurerm_network_interface.student.id,
  ]

  admin_ssh_key {
    username   = "adminusercrcr"
    public_key = tls_private_key.example_ssh.public_key_openssh
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  custom_data = base64encode(templatefile("${path.module}/tailscale_cloudinit.tpl", {
    tailscale_auth_key = var.tailscale_auth_key
  }))
}


resource "azurerm_key_vault_secret" "publicsshkey" {
  name         = "student-ssh-key-public"
  value        = tls_private_key.example_ssh.public_key_openssh
  key_vault_id = azurerm_key_vault.kvservice.id
  tags         = local.common_tags
  depends_on   = [azurerm_key_vault.kvservice]

}
resource "azurerm_key_vault_secret" "sshkey" {
  name         = "student-ssh-key-private"
  value        = tls_private_key.example_ssh.private_key_openssh
  key_vault_id = azurerm_key_vault.kvservice.id
  tags         = local.common_tags
  depends_on   = [azurerm_key_vault.kvservice]

}*/