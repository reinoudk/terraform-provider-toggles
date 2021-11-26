terraform {
  required_providers {
    toggles = {
      source = "reinoudk/toggles"
      version = "0.2.0"
    }
  }
  required_version = "~> 1.0"
}

resource "time_rotating" "toggle_interval" {
  rotation_minutes = 1
}

resource "toggles_leapfrog" "toggle" {
  # Optional, remove to toggle on each apply
  trigger = time_rotating.toggle_interval.rotation_rfc3339
}

resource "random_string" "alpha" {
  length = 10

  keepers = {
    alpha = toggles_leapfrog.toggle.alpha_timestamp
  }
}

resource "random_string" "beta" {
  length = 10

  keepers = {
    beta = toggles_leapfrog.toggle.beta_timestamp
  }
}

output "current_secret" {
  value = toggles_leapfrog.toggle.alpha ? random_string.alpha : random_string.beta
}
