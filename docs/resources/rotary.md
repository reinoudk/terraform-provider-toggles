---
page_title: "rotary Resource - terraform-provider-toggles"
subcategory: ""
description: |-
  The rotary resource allows you to change the value of n outputs in a rotating fashion.
---

# Resource `toggles_rotary`

The rotary resource allows you to change the value of n outputs in a rotating fashion. This is useful when you want to
rotate a resource but want to keep previous versions around as well. It's a more powerful version of the `leapfrog`
resource but uses counters instead of timestamps to signal changes. This is due to a limitation in the SDK, that
prevents marking a single items in a list as having a new computed value.

## Example Usage

```terraform
resource "time_rotating" "toggle_interval" {
  rotation_hours = 1
}

locals {
  n = 3
}

resource "toggles_rotary" "toggle" {
  n = local.n
  # Optional, remove to toggle on each apply
  trigger = time_rotating.toggle_interval.rotation_rfc3339
}

resource "google_service_account_key" "keys" {
  service_account_id = google_service_account.account.name

  count = local.n

  keepers = {
    rotate = toggles_rotary.toggle.counters[count.index]
  }
}

output "newest_key" {
  value = google_service_account_key[toggles_rotary.toggle.active_output]
  sensitive = true
}
```

## Argument Reference

- `trigger` - (Optional) An arbitrary string value that, when changed, toggles the output. Use this to set the min
  cadence of toggling the output. If left empty, the toggle is switched on each apply.
- `n` - The number of outputs. Should be between 2 and 256.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `outputs` - A list of n boolean outputs. The value indicates whether the output is active (was changed last).
- `active_output` - The 0-index based number of the active output.
- `counters` - A list of counters denoting the number of times the corresponding output was set to true.
