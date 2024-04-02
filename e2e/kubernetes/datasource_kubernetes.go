package kubernetes

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceKubernetesService() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Kubernetes service",
				ForceNew:    true,
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service ID of the Kubernetes Cluster",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the Kubernetes Cluster",
				ForceNew:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the block storage",
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				// ForceNew:    true,
				Description: "ID of the project. It should be unique",
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew:    true,
				Description: "Location of the block storage",
				ValidateFunc: validation.StringInSlice([]string{
					"Delhi",
					"Mumbai",
				}, false),
				Default: "Delhi",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the Kubernetes Service",
			},
			"master_node_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the master node of the kubernetes cluster",
				ForceNew:    true,
			},
		},
		ReadContext: dataSourceReadKubernetes,
	}
}

func dataSourceReadKubernetes(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("=============INSIDE KUBERNETES READ DATA SOURCE==========================")
	kubernetesId := d.Get("service_id").(string)
	location := d.Get("location").(string)
	kubernetes, err := apiClient.GetKubernetesServiceInfo(kubernetesId, location, d.Get("project_id").(int))
	if err != nil {
		return diag.Errorf("error finding Item with ID %s", kubernetesId)
	}

	log.Printf("[INFO] KUBERNETES READ | BEFORE SETTING DATA")
	data := kubernetes["data"].([]interface{})[0].(map[string]interface{})
	serviceIDFloat, ok := data["service_id"].(float64)
	if !ok {
		return diag.Errorf("Failed to convert 'service_id' to float64")
	}
	serviceIDStr := fmt.Sprintf("%.0f", serviceIDFloat)
	d.SetId(serviceIDStr)
	d.Set("name", data["service_name"].(string))
	d.Set("status", data["state"].(string))
	d.Set("version", data["version"].(string))
	d.Set("created_at", data["created_at"].(string))
	d.Set("master_node_id", strconv.FormatFloat(data["master_node_id"].(float64), 'f', -1, 64))
	return diags
}
