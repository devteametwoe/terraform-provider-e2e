package image

import (
	"context"
	// "fmt"
	"log"
	// "math"
	// "regexp"
	//"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceImages() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mention the region of the images you want to list",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "project id associated to images",
			},
			"image_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all the saved images",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "This id is used to create a node using the image",
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the image",
						},
						"os_distribution": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the image",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the image",
						},
						"distro": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "type of distro used",
						},
						"sku_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current state of the image",
						},
						"timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},

		ReadContext: dataSourceReadImages,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceReadImages(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	log.Printf("[INFO] Inside images data source ")
	Response, err := apiClient.GetSavedImages(d.Get("location").(string), d.Get("project_id").(string))
	if err != nil {
		return diag.Errorf("error finding saved images")
	}

	d.Set("image_list", flattenImages(&Response.Data))
	d.SetId("saved_image_list")
	var diags diag.Diagnostics
	return diags
}
func flattenImages(imageList *[]models.Image) []interface{} {

	if imageList != nil {

		ois := make([]interface{}, len(*imageList), len(*imageList))

		for i, image := range *imageList {

			oi := make(map[string]interface{})
			oi["template_id"] = image.Template_id
			oi["distro"] = image.Distro
			oi["image_id"] = image.Image_id
			oi["image_state"] = image.Image_state
			oi["image_type"] = image.Image_type
			oi["name"] = image.Name
			oi["sku_type"] = image.Sku_type
			oi["os_distribution"] = image.Os_distribution
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
