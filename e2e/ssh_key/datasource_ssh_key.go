package ssh_key

import (
	"context"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSshKey() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The label(name) of the ssh key",
				ForceNew:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project associated with the ssh key",
			},
			"ssh_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "your ssh key",
			},
			"project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the project associated with the ssh key",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of the SSH Key",
			},
		},

		ReadContext: dataSourceReadSshKey,
	}
}
func dataSourceReadSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	log.Printf("[INFO] INSIDE SSH KEY DATA SOURCE | read")
	label := d.Get("label").(string)
	project_id := d.Get("project_id").(string)
	res, err := apiClient.GetSshKey(label, project_id)
	if err != nil {
		return diag.Errorf("error finding ssh key with label %s | %+v", label, err)
	}

	data := res["data"].(map[string]interface{})
	log.Printf("[INFO] SSH KEY DATA SOURCE | READ | data : %+v", data)

	ssh_key_id := strconv.FormatFloat(data["pk"].(float64), 'f', 0, 64)

	d.SetId(ssh_key_id)

	d.Set("label", data["label"].(string))
	d.Set("ssh_key", data["ssh_key"].(string))
	d.Set("project_name", data["project_name"].(string))
	d.Set("timestamp", data["timestamp"].(string))
	log.Printf("[INFO] NODE DATA SOURCE | d : %+v", *d)

	return diags

}
