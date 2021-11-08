---
page_title: "leapfrog Resource - terraform-provider-toggles"
subcategory: ""
description: |-
  The leapfrog resource allows you to change a single value of two outputs in an alternating fashion. 
---

# Resource `toggles_leapfrog`

The leapfrog resource allows you to change a single value of two outputs in an alternating fashion.

## Example Usage

```terraform
resource "time_rotating" "toggle_interval" {
  rotation_hours = 1
}

resource "toggles_leapfrog" "toggle" {
  trigger = time_rotating.toggle_interval.rotation_rfc3339
}

resource "google_service_account_key" "key_blue" {
  service_account_id = google_service_account.account.name

  keepers = {
    rotation = toggle_leapfrog.toggle.blue
  }
}

resource "google_service_account_key" "key_green" {
  service_account_id = google_service_account.account.name

  keepers = {
    rotation = toggle_leapfrog.toggle.green
  }
}

output "active_key" {
  value = toggles_leapfrog.blue ? google_service_account_key.key_blue : google_service_account_key.key_green 
}
```

## Argument Reference

- `trigger` - (Optional) An arbitrary string value that, when changed, toggles the output. Use this to set the min
cadence of toggling the output.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `blue_timestamp` - An UTC RFC333 timestamp denoting the last time the blue value was updated.
- `green_timestamp` - An UTC RFC333 timestamp denoting the last time the blue value was updated.
- `blue` - A boolean indicating whether the blue output is active (changed last). This is always the inverse of green.
- `green` - A boolean indicating whether the green output is active (changed last). This is always the inverse of blue.
