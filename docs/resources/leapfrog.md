---
page_title: "leapfrog Resource - terraform-provider-toggles"
subcategory: ""
description: |-
  The leapfrog resource allows you to change a single value of two outputs in an alternating fashion. 
---

# Resource `toggles_leapfrog`

The leapfrog resource allows you to change a single value of two outputs in an alternating fashion. This is useful when
you want to rotate a resource but always keep the previous version around as well.

## Example Usage

```terraform
resource "time_rotating" "toggle_interval" {
  rotation_hours = 1
}

resource "toggles_leapfrog" "toggle" {
  # Optional, remove to toggle on each apply
  trigger = time_rotating.toggle_interval.rotation_rfc3339
}

resource "google_service_account_key" "alpha" {
  service_account_id = google_service_account.account.name

  keepers = {
    alpha = toggles_leapfrog.toggle.alpha_timestamp
  }
}

resource "google_service_account_key" "beta" {
  service_account_id = google_service_account.account.name

  keepers = {
    beta = toggles_leapfrog.toggle.beta_timestamp
  }
}

output "current_key" {
  value = toggles_leapfrog.toggle.alpha ? google_service_account_key.alpha : google_service_account_key.beta
  sensitive = true
}
```

## Argument Reference

- `trigger` - (Optional) An arbitrary string value that, when changed, toggles the output. Use this to set the min
cadence of toggling the output. If left empty, the toggle is switched on each apply.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `alpha_timestamp` - An UTC RFC333 timestamp denoting the last time the alpha value was updated.
- `beta_timestamp` - An UTC RFC333 timestamp denoting the last time the beta value was updated.
- `alpha` - A boolean indicating whether the alpha output is active (changed last). This is always the inverse of beta.
- `beta` - A boolean indicating whether the beta output is active (changed last). This is always the inverse of alpha.
