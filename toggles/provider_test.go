package toggles

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"toggles": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}

}