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
			"blue_timestamp": {
				Type: schema.TypeString,
				Description: "An UTC RFC333 timestamp denoting the last time the blue value was updated.",
				Computed: true,
			},
			"green_timestamp": {
				Type: schema.TypeString,
				Description: "An UTC RFC333 timestamp denoting the last time the green value was updated.",
				Computed: true,
			},
			"blue": {
				Type: schema.TypeBool,
				Description: "A boolean indicating whether the blue output is active (changed last). This is always the inverse of green.",
				Computed: true,
			},
			"green": {
				Type: schema.TypeBool,
				Description: "A boolean indicating whether the green output is active (changed last). This is always the inverse of blue.",
				Computed: true,
			},
		},
	}
}


// Create a new resource
func resourceLeapfrogCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("blue", true); err != nil {
		diags = diag.Errorf("could not set blue", err)
	}

	if err := d.Set("green", false); err != nil {
		diags = diag.Errorf("could not set green", err)
	}

	now := time.Now().Format(time.RFC3339)

	if err := d.Set("blue_timestamp", now); err != nil {
		diags = diag.Errorf("could not set blue_timestamp", err)
	}

	if err := d.Set("green_timestamp", now); err != nil {
		diags = diag.Errorf("could not set green_timestamp", err)
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

	blue := d.Get("blue").(bool)
	green := d.Get("green").(bool)

	blue = !blue
	green = !green

	if err := d.Set("blue", blue); err != nil {
		diags = diag.Errorf("could not set blue", err)
	}

	if err := d.Set("green", green); err != nil {
		diags = diag.Errorf("could not set green", err)
	}

	now := time.Now().Format(time.RFC3339)

	if blue {
		if err := d.Set("blue_timestamp", now); err != nil {
			diags = diag.Errorf("could not set blue_timestamp", err)
		}
	}

	if green {
		if err := d.Set("green_timestamp", now); err != nil {
			diags = diag.Errorf("could not set green_timestamp", err)
		}
	}

	return diags
}

// Delete resource
func resourceLeapfrogDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics  {
	var diags diag.Diagnostics

	return diags
}
