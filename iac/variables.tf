

variable "globals" {
  type = map(any)

  default = {
    location               = "West Europe"
  }
}
variable "network" {
    type = map(any)

    # most Austrian University VPNs
    default = { 
        allowlist_ips = "128.130.0.0/15,193.171.80.0/21,193.170.16.0/20,193.170.185.0/24,129.27.0.0/16,138.232.0.0/16,141.244.0.0/16,143.50.0.0/16,143.205.0.0/16,140.78.0.0/16,193.186.176.0/22,193.186.172.0/22,149.148.0.0/16,192.82.158.0/24,193.171.96.0/21,193.171.104.0/22"
    }
  
}

variable "myuser"{
  default = "56ea78b9-6d9f-495b-85ac-7caa86ccc191"
}