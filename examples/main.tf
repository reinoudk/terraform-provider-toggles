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
  value = toggles_leapfrog.toggle.blue ? toggles_leapfrog.toggle.blue_timestamp : toggles_leapfrog.toggle.green_timestamp
}

output "active_color" {
  value = toggles_leapfrog.toggle.blue ? "blue" : "green"
}

output "blue_timestamp" {
  value = toggles_leapfrog.toggle.blue_timestamp
}

output "green_timestamp" {
  value = toggles_leapfrog.toggle.green_timestamp
}
