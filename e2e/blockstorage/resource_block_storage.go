package blockstorage

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
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
				Description: "Status of the node",
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

	if d.HasChange("vm_id") {
		prevVMID, currVMID := d.GetChange("vm_id")
		prevName, _ := d.GetChange("name")
		prevSize, _ := d.GetChange("size")
		log.Printf("[INFO] prevVMID %v, currVMID %v, type(currVMID) %T", prevVMID, currVMID, currVMID)

		if d.Get("status") == "Attached" && prevVMID != "" {
			vm_id, err := strconv.Atoi(prevVMID.(string))
			if err != nil {
				setPrevState(d, prevVMID, prevName, prevSize)
				return diag.FromErr(err)
			}
			blockStorage := models.BlockStorageAttach{
				VM_ID: vm_id,
			}
			res, err := apiClient.AttachOrDetachBlockStorage(&blockStorage, "detach", blockStorageID, project_id, location)
			log.Printf("[INFO] BLOCK STORAGE DETACH | RESPONSE BODY | %+v", res)
			if err != nil {
				setPrevState(d, prevVMID, prevName, prevSize)
				return diag.FromErr(err)
			}
			d.Set("status", "Available")
			log.Printf("[INFO] BLOCK STORAGE DETACH | RESPONSE BODY | %+v", res)

		}
		waitForDetach(apiClient, blockStorageID, project_id, location)

		if currVMID != "" && currVMID != nil {
			if d.Get("status") == "Available" {
				vm_id, err := strconv.Atoi(currVMID.(string))
				if err != nil {
					setPrevState(d, "", prevName, prevSize)
					return diag.FromErr(err)
				}
				blockStorage := models.BlockStorageAttach{
					VM_ID: vm_id,
				}
				resBlockStorage, err := apiClient.AttachOrDetachBlockStorage(&blockStorage, "attach", blockStorageID, project_id, location)
				log.Printf("[INFO] BLOCK STORAGE ATTACH | RESPONSE BODY | %+v", resBlockStorage)
				if err != nil {
					setPrevState(d, "", prevName, prevSize)
					return diag.FromErr(err)
				}

				log.Printf("[INFO] BLOCK STORAGE ATTACH | RESPONSE BODY | %+v", resBlockStorage)
				if _, codeok := resBlockStorage["code"]; !codeok {
					setPrevState(d, "", prevName, prevSize)
					return diag.Errorf(resBlockStorage["message"].(string))
				}
				return diags
			} else {
				setPrevState(d, prevVMID, prevName, prevSize)
				return diag.Errorf("block storage cannot be attached to a node unless it is in available state")
			}
		}

	}

	if d.HasChange("size") || d.HasChange("name") {
		prevName, currName := d.GetChange("name")
		prevSize, currSize := d.GetChange("size")
		err := validateSize(d, m)
		if err != nil {
			d.Set("name", prevName)
			d.Set("size", prevSize)
			return diag.FromErr(err)
		}
		log.Printf("[INFO] prevSize %v, currSize %v", prevSize, currSize)

		if d.Get("status") == "Attached" {
			tolerance := 1e-6
			if currSize.(float64) > prevSize.(float64)+tolerance {
				log.Printf("[INFO] BLOCK STORAGE UPGRADE STARTS")
				vmID := blockStorage["data"].(map[string]interface{})["vm_detail"].(map[string]interface{})["vm_id"]
				blockStorage := models.BlockStorageUpgrade{
					Name:  currName.(string),
					Size:  currSize.(float64),
					VM_ID: vmID.(float64),
				}
				log.Printf("[INFO] BlockStorage details: %+v %T", blockStorage, blockStorage.VM_ID)
				resBlockStorage, err := apiClient.UpdateBlockStorage(&blockStorage, blockStorageID, project_id, location)
				if err != nil {
					d.Set("size", prevSize)
					d.Set("name", prevName)
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
			return diag.Errorf("You cannot change the block storage size and name unless you are upgrading it")
		} else {
			d.Set("size", prevSize)
			d.Set("name", prevName)
			return diag.Errorf("You cannot upgrade a block storage name or size unless it is attached to a node")
		}
	}
	return resourceReadBlockStorage(ctx, d, m)
}

func resourceDeleteBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	blockStorageID := d.Id()
	node_status := d.Get("status").(string)
	if node_status == "Saving" || node_status == "Creating" {
		return diag.Errorf("Node in %s state", node_status)
	}
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

func setPrevState(d *schema.ResourceData, prevVMID, prevName, prevSize interface{}) {
	d.Set("vm_id", prevVMID)
	d.Set("name", prevName)
	d.Set("size", prevSize)
}

func waitForDetach(apiClient *client.Client, blockStorageID string, project_id int, location string) diag.Diagnostics {
	for {
		blockStorage, err := apiClient.GetBlockStorage(blockStorageID, project_id, location)
		if err != nil {
			log.Printf("[ERROR] Error getting block storage %s", err)
			return diag.FromErr(err)
		}
		data := blockStorage["data"].(map[string]interface{})
		if data["status"] == "Available" {
			break
		}
		// Wait for 2 seconds before checking the status again (is Volume Detached?)
		time.Sleep(2 * time.Second)
	}
	return nil
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
