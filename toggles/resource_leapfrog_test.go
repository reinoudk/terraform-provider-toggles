package toggles

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccLeapfrog(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Applying the resource for the first time should set alpha to active and initialize both timestamps
				// with equal values.
				Config: testAccLeapfrogResource("initial"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("toggles_leapfrog.test", "alpha", "true"),
					resource.TestCheckResourceAttr("toggles_leapfrog.test", "beta", "false"),
					resource.TestCheckResourceAttrPair("toggles_leapfrog.test", "alpha_timestamp", "toggles_leapfrog.test", "beta_timestamp"),
					testAccValidRFC3339("toggles_leapfrog.test", "alpha_timestamp"),
					testAccValidRFC3339("toggles_leapfrog.test", "beta_timestamp"),
				),
			},
			{
				// Re-applying the resource with an un-changed trigger value should have the same output.
				Config: testAccLeapfrogResource("initial"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("toggles_leapfrog.test", "alpha", "true"),
					resource.TestCheckResourceAttr("toggles_leapfrog.test", "beta", "false"),
					resource.TestCheckResourceAttrPair("toggles_leapfrog.test", "alpha_timestamp", "toggles_leapfrog.test", "beta_timestamp"),
					testAccValidRFC3339("toggles_leapfrog.test", "alpha_timestamp"),
					testAccValidRFC3339("toggles_leapfrog.test", "beta_timestamp"),
				),
			},
			{
				// Re-applying the resource with a changed trigger value should mark beta as active.
				Config: testAccLeapfrogResource("change-1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("toggles_leapfrog.test", "alpha", "false"),
					resource.TestCheckResourceAttr("toggles_leapfrog.test", "beta", "true"),
					testAccTimeAfter("toggles_leapfrog.test", "beta_timestamp", "toggles_leapfrog.test", "alpha_timestamp"),
				),
			},
			{
				// Re-applying the resource with a changed trigger value again should mark alpha as active again.
				Config: testAccLeapfrogResource("change-2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("toggles_leapfrog.test", "alpha", "true"),
					resource.TestCheckResourceAttr("toggles_leapfrog.test", "beta", "false"),
					testAccTimeAfter("toggles_leapfrog.test", "alpha_timestamp", "toggles_leapfrog.test", "beta_timestamp"),
				),
			},
		},
	})
}

func testAccLeapfrogResource (trigger string) string {
	return fmt.Sprintf(`
resource "toggles_leapfrog" "test" {
  trigger = "%s"
}
`, trigger)
}
