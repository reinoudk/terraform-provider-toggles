terraform {
  required_providers {
    toggles = {
      source  = "reinoud.dev/dev/toggles"
    }
  }
  required_version = "~> 1.0"
}



resource "time_rotating" "toggle_interval" {
  rotation_minutes = 1
}

resource "toggles_leapfrog" "toggle" {
  trigger = time_rotating.toggle_interval.rotation_rfc3339
}

output "latest_timestamp" {
  value = toggles_leapfrog.toggle.alpha ? toggles_leapfrog.toggle.alpha_timestamp : toggles_leapfrog.toggle.beta_timestamp
}

output "active_color" {
  value = toggles_leapfrog.toggle.alpha ? "alpha" : "beta"
}

output "alpha_timestamp" {
  value = toggles_leapfrog.toggle.alpha_timestamp
}

output "beta_timestamp" {
  value = toggles_leapfrog.toggle.beta_timestamp
}
