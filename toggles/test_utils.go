package toggles

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"time"
)

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

// The sleep function sleeps for 1 second, to allow time to pass
func sleep() {
	time.Sleep(1 * time.Second)
}
