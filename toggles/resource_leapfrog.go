package toggles

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourceLeapfrog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLeapfrogCreate,
		ReadContext:   resourceLeapfrogRead,
		UpdateContext: resourceLeapfrogUpdate,
		DeleteContext: resourceLeapfrogDelete,
		CustomizeDiff: customizeDiffLeapfrog,
		Schema: map[string]*schema.Schema {
			"trigger": {
				Type: schema.TypeString,
				Description: "An arbitrary string value that, when changed, toggles the output.",
				Optional: true,
			},
			"alpha_timestamp": {
				Type: schema.TypeString,
				Description: "An UTC RFC333 timestamp denoting the last time the alpha value was updated.",
				Computed: true,
			},
			"beta_timestamp": {
				Type: schema.TypeString,
				Description: "An UTC RFC333 timestamp denoting the last time the beta value was updated.",
				Computed: true,
			},
			"alpha": {
				Type: schema.TypeBool,
				Description: "A boolean indicating whether the alpha output is active (changed last). This is always the inverse of beta.",
				Computed: true,
			},
			"beta": {
				Type: schema.TypeBool,
				Description: "A boolean indicating whether the beta output is active (changed last). This is always the inverse of alpha.",
				Computed: true,
			},
		},
	}
}

// customizeDiffLeapfrog ensures that we show changes in the diff phase.
// During creation it is responsive for setting the initial values of alpha and beta.
// During an update it is responsible for toggling alpha and beta, and marking the timestamps with new computed values,
// in a leapfrog fashion.
func customizeDiffLeapfrog(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// New resource: only set alpha and beta now. The timestamps are set in resourceLeapfrogCreate
	if d.Id() == "" {
		if err := d.SetNew("alpha", true); err != nil {
			return fmt.Errorf("could not set alpha: %+v", err)
		}

		if err := d.SetNew("beta", false); err != nil {
			return fmt.Errorf("could not set beta: %+v", err)
		}

		return nil
	}

	// If the trigger is set, but does not have a change, we shouldn't change anything.
	// If the trigger is empty, always update.
	trigger := d.Get("trigger").(string)
	if trigger != "" && !d.HasChange("trigger")  {
		return nil
	}

	alpha := d.Get("alpha").(bool)
	beta := d.Get("beta").(bool)

	alpha = !alpha
	beta = !beta

	if err := d.SetNew("alpha", alpha); err != nil {
		return fmt.Errorf("could not set alpha: %+v", err)
	}

	if err := d.SetNew("beta", beta); err != nil {
		return fmt.Errorf("could not set beta: %+v", err)
	}

	if alpha {
		if err := d.SetNewComputed("alpha_timestamp"); err != nil {
			return fmt.Errorf("could not mark alpha_timestamp as new computed: %+v", err)
		}
	}

	if beta {
		if err := d.SetNewComputed("beta_timestamp"); err != nil {
			return fmt.Errorf("could not mark beta_timestamp as new computed: %+v", err)
		}
	}

	return nil
}

// resourceLeapfrogCreate set the initial timestamps.
// The initial alpha and beta values are set in customizeDiffLeapfrog.
func resourceLeapfrogCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	now := time.Now().Format(time.RFC3339)

	if err := d.Set("alpha_timestamp", now); err != nil {
		diags = diag.Errorf("could not set alpha_timestamp: %+v", err)
	}

	if err := d.Set("beta_timestamp", now); err != nil {
		diags = diag.Errorf("could not set beta_timestamp: %+v", err)
	}

	// Not important
	d.SetId("toggle")

	return diags
}

// resourceLeapfrogRead is a noop as all attributes are internal.
func resourceLeapfrogRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics  {
	var diags diag.Diagnostics

	return diags
}

// resourceLeapfrogUpdate updates the timestamps depending on whether alpha or beta is active.
func resourceLeapfrogUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics  {
	var diags diag.Diagnostics

	alpha := d.Get("alpha").(bool)
	beta := d.Get("beta").(bool)

	now := time.Now().Format(time.RFC3339)

	if alpha {
		if err := d.Set("alpha_timestamp", now); err != nil {
			diags = diag.Errorf("could not set alpha_timestamp: %+v", err)
		}
	}

	if beta {
		if err := d.Set("beta_timestamp", now); err != nil {
			diags = diag.Errorf("could not set beta_timestamp: %+v", err)
		}
	}

	return diags
}

// resourceLeapfrogDelete is a noop, as no external resource are being managed.
func resourceLeapfrogDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics  {
	var diags diag.Diagnostics

	return diags
}
