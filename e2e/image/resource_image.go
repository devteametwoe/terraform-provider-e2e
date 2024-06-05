package image

import (
	// "context"

	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	//"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	// "github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceImage() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"node_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ValidateFunc: validateName,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "location of the image",
				ValidateFunc: validation.StringInSlice([]string{
					"Delhi",
					"Mumbai",
					"Delhi-NCR-2",
				}, false),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the image",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "project_id",
				ForceNew:    true,
			},
			"template_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "This id is used to create a node using the image",
			},
			"image_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current state of the image",
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
			"distro": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of distro used",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of distro used",
			},
		},

		CreateContext: resourceCreateImage,
		ReadContext:   resourceReadImage,
		UpdateContext: resourceUpdateImage,
		DeleteContext: resourceDeleteImage,
		Exists:        resourceExistsImage,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func validateName(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceCreateImage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] IMAGE CREATE")

	resImage, err := apiClient.UpdateNode(d.Get("node_id").(string), "save_images", d.Get("name").(string), d.Get("project_id").(string), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if _, codeok := resImage.(map[string]interface{})["code"]; !codeok {
		return diag.Errorf(resImage.(map[string]interface{})["message"].(string))
	}

	data := resImage.(map[string]interface{})["data"].(map[string]interface{})
	imageId := data["id"].(float64)
	fmt.Println(data)
	log.Printf("[INFO] IMAGE CREATION | before setting fields")
	d.SetId(strconv.FormatFloat(imageId, 'f', -1, 64))
	resourceReadImage(ctx, d, m)

	return diags
}

func resourceReadImage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	log.Printf("[info] inside node Resource read")
	imageId := d.Id()

	imageres, err := apiClient.GetImage(imageId, d.Get("project_id").(string))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.Errorf("error finding Item with ID %s", imageId)
		}
	}
	log.Printf("[info] IMAGE READ | BEFORE SETTING DATA %+v", imageres)
	data := imageres.Data
	d.Set("image_state", data.Image_state)
	d.Set("template_id", data.Template_id)
	d.Set("image_type", data.Image_type)
	d.Set("creation_time", data.Creation_time)
	d.Set("os_distribution", data.Os_distribution)
	d.Set("distro", data.Distro)

	return diags

}

func resourceUpdateImage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// apiClient := m.(*client.Client)

	// nodeId := d.Id()

	// _, err := apiClient.GetNode(nodeId)
	// if err != nil {

	// 	return diag.Errorf("error finding Item with ID %s", nodeId)

	// }

	return resourceReadImage(ctx, d, m)

}

func resourceDeleteImage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	log.Printf("[INFO] DELETE IMAGE")
	imageId := d.Id()

	err := apiClient.DeleteImage(imageId, d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceExistsImage(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	ImageId := d.Id()
	_, err := apiClient.GetImage(ImageId, d.Get("project_id").(string))

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
