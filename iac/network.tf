# a first network : vnets are larger network grouping, typically segmented via CIDR ranges
# high level network design should be done on paper and widely communicated 

/*
resource "azurerm_virtual_network" "vnet-student" {
  name                = "vnet-student"
  address_space       = ["10.1.0.0/16"]
  location            = var.globals.location
  resource_group_name = data.azurerm_resource_group.rgruntime.name
}*/

# a second network 
# resource "azurerm_virtual_network" "vnet-platform" {
#   name                = "vnet-platform"
#   address_space       = ["10.2.0.0/16"]
#   location            = var.globals.location
#   resource_group_name = data.azurerm_resource_group.rgruntime.name
# }

# subnets
/*
resource "azurerm_subnet" "snet-student-vm" {
  name                 = "snet-student-vm"
  resource_group_name  = data.azurerm_resource_group.rgruntime.name
  virtual_network_name = "vnet-student"
  depends_on           = [ azurerm_virtual_network.vnet-student ]
  address_prefixes       = ["10.1.1.0/24"]
}*/
# # subnets
# # A second subnet to demo peering the two so they can talk to each other
# resource "azurerm_subnet" "snet-buildagents-k8s" {
#   name                 = "snet-buildagents-k8s"
#   resource_group_name  = data.azurerm_resource_group.rgruntime.name
#   virtual_network_name = "vnet-buildagents"
#   depends_on           = [ azurerm_virtual_network.vnet-buildagents ]
#   address_prefixes       = ["10.1.2.0/24"]
# }

# security groups need to be attached
# in Azure, they attach to subnets  (in Openstack they attach to VMs and layer2 ports)
/*
resource "azurerm_subnet_network_security_group_association" "student-vm" {
  subnet_id                 = azurerm_subnet.snet-student-vm.id
  network_security_group_id = azurerm_network_security_group.student.id
}*/

######################################## ROLE ASSIGNEMENTS ##################################
# In general, if we want terraform to assign roles to itself or other objects -> we need to be careful
# this can easily be used for privilegdge escalation
# on the other hand: if you can protect your buildagents well as well as make sure the IaC branches are safe
# it is best practise to let automation assign or generally handle as much of these fine grained settings as possible
# it is easier to not make mistakes if automation handles this for you
# 

# # give the current terraform agent the permission to Read (in order to read the network config)
# resource "azurerm_role_assignment" "networkread" {
#   scope                            = "/subscriptions/${var.subscription_id}/resourceGroups/rg-service-not2day"
#   role_definition_name             = "Reader"
#   principal_id                     =  data.azurerm_client_config.current.object_id
#   depends_on                       = [ azurerm_virtual_network.vnet-buildagents ]
# }

#resource "azurerm_role_assignment" "vmcontrib" {
#  scope                            = "/subscriptions/${var.subscription_id}/resourceGroups/rg-service-not2day"
#  role_definition_name             = "Virtual Machine Contributor"
#  principal_id                     =  data.azurerm_client_config.current.object_id
#  depends_on                       = [ azurerm_virtual_network.vnet-buildagents ]
#}