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
  trigger = time_rotating.toggle_interval.rotation_rfc3339
}

resource "google_service_account_key" "key_alpha" {
  service_account_id = google_service_account.account.name

  keepers = {
    rotation = toggle_leapfrog.toggle.alpha_timestamp
  }
}

resource "google_service_account_key" "key_beta" {
  service_account_id = google_service_account.account.name

  keepers = {
    rotation = toggle_leapfrog.toggle.beta_timestamp
  }
}

output "newest_key" {
  value = toggles_leapfrog.alpha ? google_service_account_key.key_alpha : google_service_account_key.key_beta 
}
```

## Argument Reference

- `trigger` - (Optional) An arbitrary string value that, when changed, toggles the output. Use this to set the min
cadence of toggling the output.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `alpha_timestamp` - An UTC RFC333 timestamp denoting the last time the alpha value was updated.
- `beta_timestamp` - An UTC RFC333 timestamp denoting the last time the beta value was updated.
- `alpha` - A boolean indicating whether the alpha output is active (changed last). This is always the inverse of beta.
- `beta` - A boolean indicating whether the beta output is active (changed last). This is always the inverse of alpha.
