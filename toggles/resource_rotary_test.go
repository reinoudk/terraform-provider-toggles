package toggles

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccRotary(t *testing.T) {
	n := 4

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Applying the resource for the first time should set the 0th output to active and initialize all
				// counters with equal values.
				PreConfig: sleep,
				Config: testAccRotaryResource("initial", n),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.0", "true"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.1", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.2", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.3", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.0", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.1", "0"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.2", "0"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.3", "0"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "active_output", "0"),
				),
			},
			{
				// Re-applying the resource with an un-changed trigger value should have the same output.
				PreConfig: sleep,
				Config: testAccRotaryResource("initial", n),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.0", "true"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.1", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.2", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.3", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.0", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.1", "0"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.2", "0"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.3", "0"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "active_output", "0"),
				),
			},
			{
				// Re-applying the resource with a changed trigger value should mark the old output as inactive, and the
				// next output as active. It should also increment the counter for the new active output.
				PreConfig: sleep,
				Config: testAccRotaryResource("active-1", n),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.0", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.1", "true"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.2", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.3", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.0", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.1", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.2", "0"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.3", "0"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "active_output", "1"),
				),
			},
			{
				// Re-applying the resource with a changed trigger value should mark the old output as inactive, and the
				// next output as active.
				PreConfig: sleep,
				Config: testAccRotaryResource("active-2", n),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.0", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.1", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.2", "true"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.3", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.0", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.1", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.2", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.3", "0"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "active_output", "2"),
				),
			},
			{
				// Re-applying the resource with a changed trigger value should mark the old output as inactive, and the
				// next output as active.
				PreConfig: sleep,
				Config: testAccRotaryResource("active-3", n),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.0", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.1", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.2", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.3", "true"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.0", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.1", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.2", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.3", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "active_output", "3"),
				),
			},
			{
				// Re-applying the resource with a changed trigger value should set the 0-th output to active again when
				// we wrap-around.
				PreConfig: sleep,
				Config: testAccRotaryResource("wrap-around-active-0", n),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.0", "true"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.1", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.2", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "outputs.3", "false"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.0", "2"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.1", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.2", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "counters.3", "1"),
					resource.TestCheckResourceAttr("toggles_rotary.test", "active_output", "0"),
				),
			},
		},
	})
}

func testAccRotaryResource (trigger string, n int) string {
	return fmt.Sprintf(`
resource "toggles_rotary" "test" {
  trigger = "%s"
  n = "%d"
}
`, trigger, n)
}
