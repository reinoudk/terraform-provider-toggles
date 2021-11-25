package toggles

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
	"time"
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

// testAccValidRFC3339 checks if the attribute value on the resource is a valid RFC3339 timestamp
func testAccValidRFC3339(name, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource not found: %s", name)
		}

		value, ok := rs.Primary.Attributes[attr]
		if !ok {
			return fmt.Errorf("Attribute not found: %s", name)
		}

		if _, err := time.Parse(time.RFC3339, value); err != nil {
			return fmt.Errorf("Not a valid RFC3339 timestamp: %s", value)
		}

		return nil
	}
}

// testAccTimeAfter checks whether the RFC3339 timestamp from the first attribute is after that of the second attribute
func testAccTimeAfter(nameFirst, attrFirst, nameSecond, attrSecond string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rsFirst, ok := s.RootModule().Resources[nameFirst]
		if !ok {
			return fmt.Errorf("Resource not found: %s", nameFirst)
		}

		valueFirst, ok := rsFirst.Primary.Attributes[attrFirst]
		if !ok {
			return fmt.Errorf("Attribute not found: %s", nameFirst)
		}

		timeFirst, err := time.Parse(time.RFC3339, valueFirst);
		if err != nil {
			return fmt.Errorf("Not a valid RFC3339 timestamp: %s", valueFirst)
		}

		rsSecond, ok := s.RootModule().Resources[nameSecond]
		if !ok {
			return fmt.Errorf("Resource not found: %s", nameSecond)
		}

		valueSecond, ok := rsSecond.Primary.Attributes[attrSecond]
		if !ok {
			return fmt.Errorf("Attribute not found: %s", nameSecond)
		}

		timeSecond, err := time.Parse(time.RFC3339, valueSecond);
		if err != nil {
			return fmt.Errorf("Not a valid RFC3339 timestamp: %s", valueSecond)
		}

		if !timeFirst.After(timeSecond) {
			return fmt.Errorf("Timestamp %s is not after timestamp %s", valueFirst, valueSecond)
		}

		return nil
	}
}

func testAccShouldErr(f resource.TestCheckFunc) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		if err := f(state); err != nil {
			return nil
		}
		return fmt.Errorf("Expected an error")
	}
}
