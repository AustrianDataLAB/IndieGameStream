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