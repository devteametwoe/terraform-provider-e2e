package blockstorage

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceBlockStorage() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"block_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the node to be specified to read that particular node",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource, also acts as it's unique ID",
			},
			"size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Size of the block storage in GB",
			},
			"iops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IOPS of the block storage",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the block storage",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the project. It should be unique",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Location of the block storage",
				ValidateFunc: validation.StringInSlice([]string{
					"Delhi",
					"Mumbai",
				}, false),
				Default: "Delhi",
			},
			// "created_on": {
			// 	Type:        schema.TypeString,
			// 	Computed:    true,
			// 	Description: "Creation time of the block storage",
			// },
		},
		ReadContext: dataSourceReadNode,
	}
}
func dataSourceReadNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	log.Printf("[INFO] INSIDE NODE DATA SOURCE | read")
	blockStorageID := d.Get("block_id").(string)

	blockStorage, err := apiClient.GetBlockStorage(blockStorageID, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		return diag.Errorf("error finding Item with ID %s", blockStorageID)
	}

	data := blockStorage["data"].(map[string]interface{})
	d.SetId(blockStorageID)
	log.Printf("[INFO] BLOCK STORAGE DATA SOURCE | READ | data : %+v", data)
	template := data["template"].(map[string]interface{})
	resSize := convertIntoGB(data["size"].(float64))
	d.Set("size", resSize)
	d.Set("name", data["name"].(string))
	d.Set("status", data["status"].(string))
	d.Set("iops", template["TOTAL_IOPS_SEC"].(string))
	log.Printf("[INFO] NODE DATA SOURCE | d : %+v", *d)

	return diags

}
