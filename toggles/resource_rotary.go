package toggles

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRotary() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRotaryCreate,
		ReadContext:   resourceRotaryRead,
		UpdateContext: resourceRotaryUpdate,
		DeleteContext: resourceRotaryDelete,
		CustomizeDiff: customizeDiffRotary,
		Schema: map[string]*schema.Schema {
			"trigger": {
				Type: schema.TypeString,
				Description: "An arbitrary string value that, when changed, toggles the output.",
				Optional: true,
			},
			"n": {
				Type: schema.TypeInt,
				Description: "The number of outputs. Should be between 2 and 256",
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.IntBetween(2, 256),
			},
			"outputs": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
				Description: "A list of n boolean outputs.",
				Computed: true,
			},
			"active_output": {
				Type: schema.TypeInt,
				Description: "The 0-index based number of the active output.",
				Computed: true,

			},
			"counters": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "A list of counters denoting the number of times the corresponding output was set to true.",
				Computed: true,
			},
		},
	}
}

// customizeDiffRotary ensures that we show changes in the diff phase.
// As most attributes are set during the diff-phase it functions as both the create and update function for most things.
func customizeDiffRotary(ctx context.Context, d *schema.ResourceDiff, i interface{}) error {
	// New resource: only set outputs and active_output now. The counters are set in resourceRotaryCreate
	if d.Id() == "" {
		n := d.Get("n").(int)

		outputs := make([]interface{}, n, n)
		outputs[0] = true

		if err := d.SetNew("outputs", outputs); err != nil {
			return fmt.Errorf("could not set outputs: %+v", err)
		}

		if err := d.SetNew("active_output", 0); err != nil {
			return fmt.Errorf("could not set active_output: %+v", err)
		}

		counters := make([]interface{}, n, n)
		counters[0] = 1
		//for i := 0; i < n; i++ {
		//	counters[i] = 0
		//}

		if err := d.SetNew("counters", counters); err != nil {
			return fmt.Errorf("could not set counters: %+v", err)
		}

		return nil
	}

	// If the trigger is set, but does not have a change, we shouldn't change anything.
	// If the trigger is empty, always update.
	trigger := d.Get("trigger").(string)
	if trigger != "" && !d.HasChange("trigger")  {
		return nil
	}

	n := d.Get("n").(int)

	currentActiveOutput := d.Get("active_output").(int)
	nextActiveOutput := (currentActiveOutput + 1) % n

	outputs := d.Get("outputs").([]interface{})
	outputs[currentActiveOutput] = false
	outputs[nextActiveOutput] = true

	if err := d.SetNew("outputs", outputs); err != nil {
		return fmt.Errorf("could not set outputs: %+v", err)
	}

	if err := d.SetNew("active_output", nextActiveOutput); err != nil {
		return fmt.Errorf("could not set active_output: %+v", err)
	}

	counters := d.Get("counters").([]interface{})
	count := counters[nextActiveOutput].(int)
	count += 1
	counters[nextActiveOutput] = count

	if err := d.SetNew("counters", counters); err != nil {
		return fmt.Errorf("could not set counters: %+v", err)
	}

	return nil
}

// resourceRotaryCreate ensure the resource's id is set
// The initial attribute values are set in customizeDiffRotary.
func resourceRotaryCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Not important
	d.SetId("toggle")

	return diags
}

// resourceRotaryRead is a noop as all attributes are internal.
func resourceRotaryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics  {
	var diags diag.Diagnostics

	return diags
}

// resourceRotaryUpdate is a noop because all the updates happen in customizeDiffRotary
func resourceRotaryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics  {
	var diags diag.Diagnostics

	return diags
}

// resourceRotaryDelete is a noop, as no external resource are being managed.
func resourceRotaryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics  {
	var diags diag.Diagnostics

	return diags
}
