package objectstore

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceObjectStore() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the bucket, also act as it's unique ID.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The My-Account Project where the bucket will be created.",
				ForceNew:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Region the bucket will be created",
				Default:     "Delhi",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
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
			"enabling_versioning": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable versioning for this bucket.",
			},
		},
		CreateContext: resourceCreateBucket,
		ReadContext:   resourceReadBucket,
		UpdateContext: resourceUpdateBucket,
		DeleteContext: resourceDeleteBucket,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateBucket(ctx context.Context, resourceData *schema.ResourceData, clientInterface interface{}) diag.Diagnostics {
	apiClient := clientInterface.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] BUCKET CREATE STARTS ")
	bucket := models.ObjectStorePayload{
		BucketName: resourceData.Get("name").(string),
		Region:     resourceData.Get("region").(string),
		ProjectID:  resourceData.Get("project_id").(int),
	}

	resbucket, err := apiClient.CreateBucket(&bucket)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] BUCKET CREATE | RESPONSE BODY | %+v", resbucket)
	if _, codeok := resbucket["code"]; !codeok {
		return diag.Errorf(resbucket["message"].(string))
	}

	data := resbucket["data"].(map[string]interface{})
	if data["is_credit_sufficient"] == false {
		return diag.Errorf(resbucket["message"].(string))
	}
	log.Printf("[INFO] Bucket creation | before setting fields")
	bucketId := data["id"].(float64)
	bucketId = math.Round(bucketId)
	resourceData.SetId(strconv.Itoa(int(math.Round(bucketId))))
	resourceData.Set("created_on", data["created_at"].(string))
	resourceData.Set("status", data["status"].(string))
	resourceData.Set("versioning_status", data["versioning_status"].(string))
	resourceData.Set("lifecycle_configuration_status", data["lifecycle_configuration_status"].(string))
	resourceData.Set("enabling_versioning", false)
	return diags
}

func resourceReadBucket(ctx context.Context, resourceData *schema.ResourceData, clientInterface interface{}) diag.Diagnostics {

	apiClient := clientInterface.(*client.Client)
	var diags diag.Diagnostics
	log.Printf("[info] inside node Resource read")
	bucketName := resourceData.Get("name").(string)
	location := resourceData.Get("region").(string)
	projectID := fmt.Sprint(resourceData.Get("project_id").(int))

	bucket, err := apiClient.GetBucket(bucketName, location, projectID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resourceData.SetId("")
		} else {
			return diag.Errorf("error finding Item with Name %s", bucketName)

		}
	}
	log.Printf("[info] Object Store Resource read | before setting data")
	data := bucket["data"].(map[string]interface{})
	log.Printf("[INFO] Object Store Data: %s", data)
	resourceData.Set("created_on", data["created_at"].(string))
	resourceData.Set("status", data["status"].(string))
	resourceData.Set("versioning_status", data["versioning_status"].(string))
	resourceData.Set("lifecycle_configuration_status", data["lifecycle_configuration_status"].(string))

	log.Printf("[info] Object Store Resource read | after setting data")
	if resourceData.Get("status").(string) == "Running" {
		resourceData.Set("enabling_versioning", true)
	}

	return diags

}

func resourceUpdateBucket(ctx context.Context, resourceData *schema.ResourceData, clientInterface interface{}) diag.Diagnostics {

	apiClient := clientInterface.(*client.Client)

	if resourceData.HasChange("name") {
		oldName, _ := resourceData.GetChange("name")
		resourceData.Set("name", oldName)
		return diag.Errorf("cannot change the bucket name of an object storage after creation")
	}
	bucketName := resourceData.Get("name").(string)
	projectID := fmt.Sprint(resourceData.Get("project_id").(int))
	region := resourceData.Get("region").(string)

	if resourceData.HasChange("enabling_versioning") {
		bucketstatus := resourceData.Get("status").(string)
		log.Printf("[INFO] %s ", bucketstatus)
		var action string
		if resourceData.Get("enabling_versioning").(bool) {
			action = "Enabled"
		} else {
			action = "Disabled"
		}
		resbucket, err := apiClient.SetBucketVersioning(bucketName, region, projectID, action)
		data := resbucket["data"].(map[string]interface{})
		if err != nil {
			return diag.FromErr(err)
		}
		resourceData.Set("versioning_status", data["bucket_versioning_status"].(string))
		resourceData.Set("enabling_versioning", resourceData.Get("enabling_versioning").(bool))
	}
	return resourceReadBucket(ctx, resourceData, clientInterface)

}

func resourceDeleteBucket(ctx context.Context, resourceData *schema.ResourceData, clientInterface interface{}) diag.Diagnostics {
	apiClient := clientInterface.(*client.Client)
	var diags diag.Diagnostics
	bucketName := resourceData.Get("name").(string)
	projectID := fmt.Sprint(resourceData.Get("project_id").(int))
	region := resourceData.Get("region").(string)

	err := apiClient.DeleteBucket(bucketName, region, projectID)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceData.SetId("")
	return diags
}

func resourceExistsObjectStore(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	bucketName := d.Get("name").(string)
	projectID := fmt.Sprint(d.Get("project_id").(int))
	region := d.Get("region").(string)
	_, err := apiClient.GetBucket(bucketName, projectID, region)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
