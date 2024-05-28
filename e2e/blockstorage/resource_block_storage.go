package blockstorage

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceBlockStorage() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the block storage, also acts as its unique ID",
				ValidateFunc: node.ValidateName,
			},
			"size": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "Size of the block storage in GB",
			},
			"iops": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IOPS of the block storage",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the project. It should be unique",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location of the block storage",
				ValidateFunc: validation.StringInSlice([]string{
					"Delhi",
					"Mumbai",
					"Delhi-NCR-2",
				}, false),
				Default: "Delhi",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the block storage",
			},
			"vm_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the VM to which the block storage is attached",
			},
			"vm_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the VM to which the block storage is attached",
			},
		},

		CreateContext: resourceCreateBlockStorage,
		ReadContext:   resourceReadBlockStorage,
		UpdateContext: resourceUpdateBlockStorage,
		DeleteContext: resourceDeleteBlockStorage,
		Exists:        resourceExistsBlockStorage,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	err := validateSize(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] BLOCK STORAGE CREATE STARTS ")
	blockStorage := models.BlockStorageCreate{
		Name: d.Get("name").(string),
		Size: d.Get("size").(float64),
	}

	iops := calculateIOPS(blockStorage.Size)
	blockStorage.IOPS = iops

	resBlockStorage, err := apiClient.NewBlockStorage(&blockStorage, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] BLOCK STORAGE CREATE | RESPONSE BODY | %+v", resBlockStorage)
	if _, codeok := resBlockStorage["code"]; !codeok {
		return diag.Errorf(resBlockStorage["message"].(string))
	}

	data := resBlockStorage["data"].(map[string]interface{})
	if data["is_credit_sufficient"] == false {
		return diag.Errorf(resBlockStorage["message"].(string))
	}
	log.Printf("[INFO] Block Storage creation | before setting fields")
	blockStorageIDFloat, ok := data["id"].(float64)
	if !ok {
		return diag.Errorf("Block ID is not a valid float64 in the response %v", data["id"])
	}

	blockStorageID := int(math.Round(blockStorageIDFloat))
	d.SetId(strconv.Itoa(blockStorageID))
	d.Set("iops", iops)
	return diags
}

func resourceReadBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] BLOCK STORAGE READ STARTS")
	blockStorageID := d.Id()

	blockStorage, err := apiClient.GetBlockStorage(blockStorageID, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return diags
		}
		return diag.Errorf("error finding Block Storage with ID %s: %s", blockStorageID, err.Error())
	}

	log.Printf("[INFO] BLOCK STORAGE READ | BEFORE SETTING DATA")
	data := blockStorage["data"].(map[string]interface{})
	template := data["template"].(map[string]interface{})
	vm_detail := data["vm_detail"].(map[string]interface{})
	// resSize := convertIntoGB(data["size"].(float64))
	d.Set("name", data["name"].(string))
	// d.Set("size", resSize)
	d.Set("status", data["status"].(string))
	d.Set("iops", template["TOTAL_IOPS_SEC"].(string))
	if val, ok := vm_detail["vm_id"]; ok {
		d.Set("vm_id", strconv.Itoa(int(val.(float64))))
		d.Set("vm_name", vm_detail["vm_name"].(string))
	} else {
		d.Set("vm_id", nil)
		d.Set("vm_name", nil)
	}

	log.Printf("[INFO] BLOCK STORAGE READ | AFTER SETTING DATA")

	return diags
}

func resourceUpdateBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	blockStorageID := d.Id()
	project_id := d.Get("project_id").(int)
	location := d.Get("location").(string)
	status := d.Get("status").(string)

	if status == constants.BLOCK_STORAGE_STATUS["ERROR"] {
		rollbackChanges(d)
		return diag.Errorf("Block Storage is in ERROR state.")
	}
	blockStorage, err := apiClient.GetBlockStorage(blockStorageID, project_id, location)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return diags
		}
		return diag.Errorf("error finding Block Storage with ID %s: %s", blockStorageID, err.Error())
	}
	if d.HasChange("location") {
		prevLocation, currLocation := d.GetChange("location")
		log.Printf("[INFO] prevLocation %v, currLocation %v", prevLocation, currLocation)
		d.Set("location", prevLocation)
		return diag.Errorf("Location cannot be changed once the block storage is created")
	}
	if d.HasChange("project_id") {
		prevProjectID, currProjectID := d.GetChange("project_id")
		log.Printf("[INFO] prevProjectID %v, currProjectID %v", prevProjectID, currProjectID)
		d.Set("project_id", prevProjectID)
		return diag.Errorf("Project ID cannot be changed once the block storage is created")
	}

	if d.HasChange("size") {
		prevName, currName := d.GetChange("name")
		prevSize, currSize := d.GetChange("size")
		err := validateSize(d, m)
		if err != nil {
			d.Set("name", prevName)
			d.Set("size", prevSize)
			return diag.FromErr(err)
		}
		log.Printf("[INFO] prevSize %v, currSize %v", prevSize, currSize)

		if d.Get("status") == constants.BLOCK_STORAGE_STATUS["ATTACHED"] {
			tolerance := 1e-6
			if currSize.(float64) > prevSize.(float64)+tolerance {
				log.Printf("[INFO] BLOCK STORAGE UPGRADE STARTS")
				vmID := blockStorage["data"].(map[string]interface{})["vm_detail"].(map[string]interface{})["vm_id"]
				blockStorage := models.BlockStorageUpgrade{
					Name:  currName.(string),
					Size:  currSize.(float64),
					VM_ID: vmID.(float64),
				}
				log.Printf("[INFO] BlockStorage details for update : %+v %T", blockStorage, blockStorage.VM_ID)
				// project_id_string := strconv.Itoa(project_id)
				// Err := node.WaitForDesiredState(apiClient, int(blockStorage.VM_ID), project_id_string, location)
				// if Err != nil {
				// 	d.Set("size", prevSize)
				// 	d.Set("name", prevName)
				// 	return Err
				// }

				resBlockStorage, err := apiClient.UpdateBlockStorage(&blockStorage, blockStorageID, project_id, location)
				if err != nil {
					d.Set("size", prevSize)
					d.Set("name", prevName)
					if checkErrorForSpecificMessage(err, "wrong state DISK_RESIZE_POWEROFF") {
						return diag.Errorf("%s | Currently resizing another disk on same virtual machine, Please Wait.", prevName)
					}
					return diag.FromErr(err)
				}
				log.Printf("[INFO] BLOCK STORAGE UPGRADE | RESPONSE BODY | %+v", resBlockStorage)
				if _, codeok := resBlockStorage["code"]; !codeok {
					d.Set("size", prevSize)
					d.Set("name", prevName)
					return diag.Errorf(resBlockStorage["message"].(string))
				}
				return diags
			}
			d.Set("size", prevSize)
			d.Set("name", prevName)
			return diag.Errorf("You cannot change the block storage size unless you are upgrading it")
		} else {
			d.Set("size", prevSize)
			d.Set("name", prevName)
			return diag.Errorf("You cannot upgrade a block storage size unless it is attached to a node")
		}
	}
	if !d.HasChange("size") && d.HasChange("name") {
		prevName, currName := d.GetChange("name")
		d.Set("name", prevName)
		return diag.Errorf("You cannot change the name of a blockstorage resource to %v after creation.", currName)
	}
	return resourceReadBlockStorage(ctx, d, m)
}

func resourceDeleteBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	blockStorageID := d.Id()
	status := d.Get("status").(string)
	if status == constants.BLOCK_STORAGE_STATUS["SAVING"] || status == constants.BLOCK_STORAGE_STATUS["CREATING"] {
		return diag.Errorf("Block storage in %s state", status)
	}
	if status == constants.BLOCK_STORAGE_STATUS["ATTACHED"] {
		return diag.Errorf("Block Storage is attached to a node. Detach it first")
	}
	log.Printf("[INFO] BLOCK STORAGE DELETE STARTS")
	err := apiClient.DeleteBlockStorage(blockStorageID, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceExistsBlockStorage(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	blockStorageID := d.Id()
	_, err := apiClient.GetBlockStorage(blockStorageID, d.Get("project_id").(int), d.Get("location").(string))

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func calculateIOPS(size float64) string {
	iops := size * 15
	return strconv.Itoa(int(iops))
}

func convertIntoGB(bsSizeRes float64) float64 {
	return bsSizeRes / 1024
}

func validateSize(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resPlans, err := apiClient.GetBlockStoragePlans(d.Get("project_id").(int), d.Get("location").(string))

	if err != nil {
		return err
	}
	availablePlans := resPlans["data"].([]interface{})
	isSizeValid := false

	for _, val := range availablePlans {
		plan, ok := val.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to assert val to map[string]interface{}")
		}
		size, ok := plan["bs_size"].(float64)
		// Convert TB into GB
		size = size * 1000
		if !ok {
			return fmt.Errorf("failed to assert bs_size to float64")
		}
		log.Printf("[INFO] plan Size: %v", size)
		isSizeValid = isSizeValid || (size == d.Get("size").(float64))
	}

	if !isSizeValid {
		return fmt.Errorf("BlockStorage Size not available in the plans")
	}
	return nil

}

func rollbackChanges(d *schema.ResourceData) {
	prevName, _ := d.GetChange("name")
	prevSize, _ := d.GetChange("size")

	d.Set("name", prevName)
	d.Set("size", prevSize)
}

func checkErrorForSpecificMessage(err error, message string) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), message)
}
