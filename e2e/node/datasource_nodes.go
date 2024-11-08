package node

import (
	"context"
	//"encoding/json"
	// "fmt"
	"log"
	// "math"
	// "regexp"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNodes() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region should specified",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the project associated with the node",
			},
			"nodes_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all the Nodes of your account . ",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The id of the node",
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_locked": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rescue_mode_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		ReadContext: dataSourceReadNodes,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceReadNodes(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	log.Printf("[INFO] Inside nodes data source ")
	Response, err := apiClient.GetNodes(d.Get("region").(string), d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] %v", Response)
	log.Printf("[INFO] NODES DATA SOURCE | before setting")
	d.Set("nodes_list", flattenNodes(&Response.Data))
	d.SetId("nodes_list")

	return diags
}

func flattenNodes(nodes *[]models.Node) []interface{} {

	if nodes != nil {
		ois := make([]interface{}, len(*nodes), len(*nodes))

		for i, node := range *nodes {
			oi := make(map[string]interface{})
			oi["id"] = node.ID
			oi["name"] = node.Name
			oi["is_locked"] = node.IsLocked
			oi["private_ip_address"] = node.PrivateIPAddress
			oi["public_ip_address"] = node.PublicIPAddress
			oi["rescue_mode_status"] = node.RescueModeStatus
			oi["status"] = node.Status
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
