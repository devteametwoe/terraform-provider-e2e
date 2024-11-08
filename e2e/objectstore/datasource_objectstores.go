package objectstore

import (
	"context"
	"fmt"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceObjectStores() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region should specified",
			},
			"bucket_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Buckets for the user",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Id for the My-Account bucket",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name for the My-Account bucket",
						},
						"bucket_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Size of the My-Account Bucket",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of Bucket",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Created Time of My-Account Bucket",
						},
						"versioning_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Is Versioning enabled?",
						},
						"lifecycle_configuration_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Is Lifecycle Rule Configured?",
						},
					},
				},
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Associated Project ID for buckets",
			},
		},
		ReadContext: dataSourceReadBuckets,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceReadBuckets(context context.Context, resourceDataSource *schema.ResourceData, clientInterface interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := clientInterface.(*client.Client)
	log.Printf("[INFO] ---- Execute Get Request to fetch Buckets Data. ---- ")
	project_id := fmt.Sprint(resourceDataSource.Get("project_id").(int))
	Response, err := apiClient.GetBuckets(resourceDataSource.Get("region").(string), project_id)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] BUCKETS DATA SOURCE | before setting --> %v", &Response.Data)
	eos_bucket_list := flattenBuckets(&Response.Data)
	resourceDataSource.Set("bucket_list", eos_bucket_list)
	resourceDataSource.SetId("bucket_list")
	return diags
}

func flattenBuckets(buckets *[]models.ObjectStore) []interface{} {

	if buckets != nil {
		buckets_list := make([]interface{}, len(*buckets), len(*buckets))

		for i, bucket := range *buckets {
			log.Printf("[INFO] Buckets----> %v", bucket)
			eos_bucket := make(map[string]interface{})
			eos_bucket["id"] = bucket.ID
			eos_bucket["name"] = bucket.Name
			eos_bucket["bucket_size"] = bucket.BucketSize
			eos_bucket["created_on"] = bucket.CreatedOn
			eos_bucket["status"] = bucket.Status
			eos_bucket["lifecycle_configuration_status"] = bucket.LifecycleConfigurationStatus
			eos_bucket["versioning_status"] = bucket.VersioningStatus
			buckets_list[i] = eos_bucket
		}
		return buckets_list
	}
	return make([]interface{}, 0)
}
