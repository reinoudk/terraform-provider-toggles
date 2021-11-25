terraform {
  required_providers {
    toggles = {
      source = "reinoudk/toggles"
      version = "0.3.0"
    }
  }
  required_version = "~> 1.0"
}

locals {
  n = 3
}

resource "toggles_rotary" "toggle" {
  n = local.n
}

resource "random_string" "rand" {
  length = 10
  count = local.n

  keepers = {
    rotate = toggles_rotary.toggle.counters[count.index]
  }
}

output "active_rand" {
  value = random_string.rand[toggles_rotary.toggle.active_output]
}

output "outputs" {
  value = toggles_rotary.toggle.outputs
}
