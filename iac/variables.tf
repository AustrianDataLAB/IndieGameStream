variable "globals" {
  type = map(any)

  default = {
    location               = "West Europe"
  }
}

variable "myuser"{
  default = "56ea78b9-6d9f-495b-85ac-7caa86ccc191"
}

variable "cluster_name" {
  type = string
  default = "indiegamestream-cluster"
}

variable "aks_admin_group_object_ids" {
  description = "aks admin group ids"
  type = list(string)
  default = [var.myuser, "7ab666bb-6355-4240-aa93-16bfbb9fd5f7"]
}