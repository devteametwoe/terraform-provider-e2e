package node

import (
	"context"
	// "fmt"
	"log"
	// "math"
	// "regexp"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	// "github.com/devteametwoe/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNode() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the node to be specified to read that particular node",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource, also acts as it's unique ID",
			},
			"label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the group",
			},
			"plan": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the Plan",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the node",
			},
			"memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "memory of the node",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the node",
			},
			"disk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disc info of the node",
			},
			"price": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "price details of the node",
			},
			"public_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public ip address alloted to node",
			},
			"private_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private ip address alloted to node if any",
			},
			"is_bitninja_license_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Can check if the bitninja license is active or not",
			},
			"is_locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "if the node is locked or not",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the project associated with the node",
			},
		},

		ReadContext: dataSourceReadNode,
	}
}
func dataSourceReadNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	log.Printf("[INFO] INSIDE NODE DATA SOURCE | read")
	nodeId := d.Get("node_id").(string)
	project_id := d.Get("project_id").(string)
	node, err := apiClient.GetNode(nodeId, project_id)
	if err != nil {
		return diag.Errorf("error finding Item with ID %s", nodeId)
	}

	data := node["data"].(map[string]interface{})
	d.SetId(nodeId)
	log.Printf("[INFO] NODE DATA SOURCE | READ | data : %+v", data)
	d.Set("name", data["name"].(string))
	d.Set("label", data["label"].(string))
	d.Set("plan", data["plan"].(string))
	d.Set("created_at", data["created_at"].(string))
	d.Set("memory", data["memory"].(string))
	d.Set("status", data["status"].(string))
	d.Set("disk", data["disk"].(string))
	d.Set("price", data["price"].(string))
	d.Set("is_locked", data["is_locked"].(bool))
	d.Set("public_ip_address", data["public_ip_address"].(string))
	d.Set("private_ip_address", data["private_ip_address"].(string))
	d.Set("is_bitninja_license_active", data["is_bitninja_license_active"].(bool))
	log.Printf("[INFO] NODE DATA SOURCE | d : %+v", *d)

	return diags

}
