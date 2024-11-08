package sfs

import (
	// "context"

	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	//"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSfs() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
         
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:   true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ValidateFunc: validateName,
			},
			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:   true,
				Description: "Details  of the Plan",
			},
			"vpc_id":{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:   true,
				Description: "virtual private cloud id of sfs",
			},
			"disk_size":{
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:   true,
				Description: "size of disk to be created",
			},
			"project_id":{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "size of disk to be created",
			},
			"disk_iops":{
				Type:       schema.TypeInt,
				Required:   true,
				ForceNew:   true,
				Description:  "input output per second",
			},
			"status":{
				Type:       schema.TypeString,
				Computed:   true,
				Optional:   true,
				Description:  "status will be updated after creation",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Location where node is to be launched",
				Default:     "Delhi",
			},
		
		},
		CreateContext: resourceCreateSfs,
		ReadContext:   resourceReadSfs,
		DeleteContext: resourceDeleteSfs,
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


func resourceCreateSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] NODE CREATE STARTS ")
	node := models.SfsCreate{
		Name:              d.Get("name").(string),
		Plan:              d.Get("plan").(string),
		Vpc_id:            d.Get("vpc_id").(string),
		Disk_size:         d.Get("disk_size").(int),
		Disk_iops:         d.Get("disk_iops").(int),
		
	}
	project_id:=d.Get("project_id").(string)
	location:=d.Get("region").(string)
	res_Sfs, err := apiClient.NewSfs(&node, project_id, location)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] NODE CREATE | RESPONSE BODY | %+v", res_Sfs)
	if _, codeok := res_Sfs["code"]; !codeok {
		return diag.Errorf(res_Sfs["message"].(string))
	}

	data := res_Sfs["data"].(map[string]interface{})
	if data["is_credit_sufficient"] == false {
		return diag.Errorf(res_Sfs["message"].(string))
	}
	log.Printf("[INFO] sfs creation | before setting fields")
	sfsId, ok := data["efs_id"].(float64)
	if !ok {
		return diag.Errorf("unable to retrieve valid 'id' from response")
	}
	
	d.SetId(strconv.Itoa(int(math.Round(sfsId))))

	return diags
	}

func resourceReadSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	log.Printf("[info] inside node Resource read")
	Sfs_id := d.Id()
	project_id:=d.Get("project_id").(string)
	location:=d.Get("region").(string)
	Sfs, err := apiClient.GetSfs(Sfs_id, project_id, location)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.Errorf("error finding Item with ID %s", Sfs_id)

		}
	}
	log.Printf("[info] node Resource read | before setting data")
	data := Sfs["data"].(map[string]interface{})
	d.Set("name", data["name"].(string))
	log.Printf("[info] node Resource read | after setting data")
	if d.Get("status").(string) == "Available" {
		d.Set("status", "power_on")
	}
	
	return diags

}

func resourceDeleteSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	Sfs_id := d.Id()
	project_id:=d.Get("project_id").(string)
	node_status := d.Get("status").(string)
	if node_status == "Creating" {
		return diag.Errorf("Node in %s state", node_status)
	}
	location:=d.Get("region").(string)
	err := apiClient.DeleteSFs(Sfs_id, project_id, location)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}



