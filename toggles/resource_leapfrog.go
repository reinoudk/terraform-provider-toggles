package toggles

import (
	"context"
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


// Create a new resource
func resourceLeapfrogCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("alpha", true); err != nil {
		diags = diag.Errorf("could not set alpha: %+v", err)
	}

	if err := d.Set("beta", false); err != nil {
		diags = diag.Errorf("could not set beta: %+v", err)
	}

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

// Read resource information
func resourceLeapfrogRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics  {
	var diags diag.Diagnostics

	return diags
}

// Update resource
func resourceLeapfrogUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics  {
	var diags diag.Diagnostics

	// When the trigger is set, but does not have a change, we shouldn't change anything.
	trigger := d.Get("trigger").(string)
	if trigger != "" && !d.HasChange("trigger")  {
		return diags
	}

	alpha := d.Get("alpha").(bool)
	beta := d.Get("beta").(bool)

	alpha = !alpha
	beta = !beta

	if err := d.Set("alpha", alpha); err != nil {
		diags = diag.Errorf("could not set alpha: %+v", err)
	}

	if err := d.Set("beta", beta); err != nil {
		diags = diag.Errorf("could not set beta: %+v", err)
	}

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

// Delete resource
func resourceLeapfrogDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics  {
	var diags diag.Diagnostics

	return diags
}
