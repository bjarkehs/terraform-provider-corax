# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    corax = {
      source = "registry.terraform.io/trifork/corax"
    }
  }
}

provider "corax" {
}
