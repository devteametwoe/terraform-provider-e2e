package reserve_ip

import (
	"context"
	// "fmt"
	"log"
	// "math"
	// "regexp"

	// "strconv"
	//"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceReserveIps() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "need to specify the region  (Mumbai/Delhi)",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "specify the project id",
			},
			"reserve_ips_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all the reserved ips",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reserve_id": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "id of the reserve_ip",
						},
						"appliance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of infra to which the node is attached",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ip address",
						},
						"reserved_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the type of ip address that is reserved",
						},
						"vm_id": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The id of the image",
						},
						"bought_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time at which the ip was bought ",
						},
						"vm_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the node to which the ip is attached if any",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "wheather the ip is attached or available",
						},
					},
				},
			},
		},

		ReadContext: dataSourceReadReserveIps,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceReadReserveIps(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	log.Printf("[INFO] Inside images data source ")
	Response, err := apiClient.GetReservedIps(d.Get("project_id").(string), d.Get("region").(string))
	if err != nil {
		return diag.Errorf("error finding saved images")
	}

	d.Set("reserve_ips_list", flattenReserveIps(&Response.Data))
	d.SetId("reserve_ips_list")
	var diags diag.Diagnostics
	return diags
}
func flattenReserveIps(ReserveIps *[]models.ReserveIp) []interface{} {

	if ReserveIps != nil {

		ois := make([]interface{}, len(*ReserveIps), len(*ReserveIps))
		for i, reserveip := range *ReserveIps {

			oi := make(map[string]interface{})
			oi["reserve_id"] = reserveip.ReserveID
			oi["appliance_type"] = reserveip.ApplianceType
			oi["bought_at"] = reserveip.BoughtAt
			oi["ip_address"] = reserveip.IPAddress
			oi["reserved_type"] = reserveip.ReservedType
			oi["status"] = reserveip.Status
			oi["vm_id"] = reserveip.VMID
			oi["vm_name"] = reserveip.VMName
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
