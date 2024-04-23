package node

import (
	// "context"

	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"time"
	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/constants"

	// "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/security_group"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceNode() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ValidateFunc: ValidateName,
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the group",
				Default:     "default",
			},
			"plan": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the Plan",
				ValidateFunc: ValidatePlanName,
			},
			"backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Tells you the state of your backups",
				Default:     false,
			},

			"image": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the image you have selected format :- ( os-version )",
				ValidateFunc: ValidateBlank,
			},
			"default_public_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Tells us the state of default public ip",
				Default:     false,
			},
			"disable_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "can disable password as per requirement",
				Default:     false,
			},
			"enable_bitninja": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "enable bitnija as per requirement",
				Default:     false,
			},
			"is_ipv6_availed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
				Default:     false,
			},
			"is_saved_image": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "used when Creating node from a saved image",
				Default:     false,
			},
			"start_scripts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"reserve_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Attach reserve ip as per requirement",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Vpc id as per requirement",
			},
			"saved_image_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "template id  is required when you save the node from saved images.Give the template id of the saved image. Required when is_saved_image field is true",
				Default:     nil,
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify the security group. Checkout security_groups datasource listing security groups",
				Elem: &schema.Schema{
					Type:        schema.TypeInt,
					Description: "ID of the security group",
				},
			},
			"default_sg": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Default Security Group",
			},
			"ssh_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify the label of ssh keys if required. Checkout ssh_keys datasource for listing ssh keys",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the node",
			},
			"memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Amount of RAM assigned to the node",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the node",
			},
			"disk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disc info of the node",
			},
			"price": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "price details of the node",
			},
			"public_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public ip address alloted to node",
			},
			"private_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private ip address alloted to node if any",
			},
			"is_bitninja_license_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Can check if the bitninja license is active or not",
			},
			"power_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "power_on",
				Description: "power_on to start the node and power_off to power off the node",
			},
			"lock_node": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Node is locked when set true .Can specify wheather to lock the node or not",
			},
			"reboot_node": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When set true node will be rebooted. Node should be in running state to perform rebooting.Alaways check the field. If you have an active disk-intensive process such as database, backups running, then a rebooting may lead to data corruption and data loss (best option is to reboot the machine from within Operating System). ",
			},
			"reinstall_node": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "for reinstalling the node. Node should be in running state to perform this action. Always check this field as it will delete all your data permenantly when set true.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the project associated with the node",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Delhi",
				Description: "Location where you want to create node.(ex - \"Delhi\", \"Mumbai\").",
				ValidateFunc: validation.StringInSlice([]string{
					"Delhi",
					"Mumbai",
					"Delhi-NCR-2",
				}, false),
			},
			"vm_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The id of the VM.",
			},
			"block_storage_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The id of the block storage to be attached to the node",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					Description:  "ID of the block storage",
					ValidateFunc: ValidateBlank,
				},
			},
		},

		CreateContext: resourceCreateNode,
		ReadContext:   resourceReadNode,
		UpdateContext: resourceUpdateNode,
		DeleteContext: resourceDeleteNode,
		Exists:        resourceExistsNode,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func ValidateName(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	if len(value) == 0 {
		errs = append(errs, fmt.Errorf("name cannot be empty"))
		return warns, errs
	}
	validNameRegexp := regexp.MustCompile(`^[a-zA-Z0-9-_]{1,50}$`)
	if !validNameRegexp.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("the name field cannot be blank, must not contain whitespace or special characters, and must be between 1 and 50 characters in length. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceCreateNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	copy_ssh_keys := d.Get("ssh_keys")
	new_SSH_keys, Err := convertLabelToSshKey(m, d.Get("ssh_keys").([]interface{}), d.Get("project_id").(string))

	if Err != nil {
		return Err
	}
	d.Set("ssh_keys", new_SSH_keys)

	if len(d.Get("block_storage_ids").([]interface{})) > 1 {
		return diag.Errorf("Can only attach a single block storage while node creation.")
	}
	image_id := 0
	if len(d.Get("block_storage_ids").([]interface{})) == 1 {
		if d.Get("plan").(string)[0:2] == "C2" {
			return diag.Errorf("Block storage can not be attached to C2 plan")
		}
		image_id_temp, err := convertStringToInt(d.Get("block_storage_ids").([]interface{})[0].(string))
		if err != nil {
			return diag.FromErr(err)
		}
		image_id = image_id_temp
	}

	log.Printf("[INFO] NODE CREATE STARTS ")
	response, err := apiClient.GetSecurityGroupList(d.Get("project_id").(string), d.Get("location").(string))
	log.Printf("[INFO] GET Security groups | RESPONSE BODY | %+v", response)
	if err != nil {
		log.Printf("[ERROR] Error getting Security Group List inside Node Create. Error : %s", err)
		return diag.Errorf("please confirm the project_id or location that you defined.")
	}
	defaultSG := getDefaultSG(response)
	d.Set("default_sg", defaultSG)

	security_group := defaultSG
	if securityGroupsList, ok := d.GetOk("security_group_ids"); ok {
		if securityGroupsList != nil {
			if securityGroups, ok := securityGroupsList.([]interface{}); ok && len(securityGroups) > 0 {
				security_group = securityGroups[0].(int)
				if len(securityGroups) > 1 {
					log.Printf("Can only attach a single security group while node creation. Only the first Security Group will be attached")
					d.Set("security_group_ids", []int{security_group})
				}
			}
		}
	}

	node := models.NodeCreate{
		Name:              d.Get("name").(string),
		Label:             d.Get("label").(string),
		Plan:              d.Get("plan").(string),
		Backup:            d.Get("backup").(bool),
		Image:             d.Get("image").(string),
		Default_public_ip: d.Get("default_public_ip").(bool),
		Disable_password:  d.Get("disable_password").(bool),
		Enable_bitninja:   d.Get("enable_bitninja").(bool),
		Is_ipv6_availed:   d.Get("is_ipv6_availed").(bool),
		Is_saved_image:    d.Get("is_saved_image").(bool),
		// Region:            d.Get("location").(string),
		Reserve_ip:        d.Get("reserve_ip").(string),
		Vpc_id:            d.Get("vpc_id").(string),
		Security_group_id: security_group,
		SSH_keys:          d.Get("ssh_keys").([]interface{}),
		Start_scripts:     d.Get("start_scripts").([]interface{}),
		Image_id:          image_id,
	}

	if node.Vpc_id != "" {
		vpc_details, err := apiClient.GetVpc(node.Vpc_id, d.Get("project_id").(string), d.Get("location").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		data := vpc_details.Data
		if data.State != "Active" {
			return diag.Errorf("Can not create node resource, vpc is in %s state", data.State)
		}
	}
	project_id := d.Get("project_id").(string)
	resnode, err := apiClient.NewNode(&node, project_id, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] NODE CREATE | RESPONSE BODY | %+v", resnode)
	if _, codeok := resnode["code"]; !codeok {
		return diag.Errorf(resnode["message"].(string))
	}

	data := resnode["data"].(map[string]interface{})
	if data["is_credit_sufficient"] == false {
		return diag.Errorf(resnode["message"].(string))
	}
	log.Printf("[INFO] node creation | before setting fields")
	nodeId := data["id"].(float64)
	nodeId = math.Round(nodeId)
	d.SetId(strconv.Itoa(int(math.Round(nodeId))))
	d.Set("ssh_keys", copy_ssh_keys)
	d.Set("is_active", data["is_active"].(bool))
	d.Set("created_at", data["created_at"].(string))
	d.Set("memory", data["memory"].(string))
	d.Set("status", data["status"].(string))
	d.Set("disk", data["disk"].(string))
	d.Set("price", data["price"].(string))
	d.Set("vm_id", int(data["vm_id"].(float64)))
	return diags
}

func resourceReadNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	copy_ssh_keys := d.Get("ssh_keys")
	log.Printf("[info] inside node Resource read")
	nodeId := d.Id()
	project_id := d.Get("project_id").(string)
	node, err := apiClient.GetNode(nodeId, project_id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.Errorf("error finding Item with ID %s", nodeId)

		}
	}
	log.Printf("[info] node Resource read | before setting data")
	data := node["data"].(map[string]interface{})

	d.Set("name", data["name"].(string))
	d.Set("label", data["label"].(string))
	d.Set("plan", data["plan"].(string))
	d.Set("created_at", data["created_at"].(string))
	d.Set("memory", data["memory"].(string))
	d.Set("status", data["status"].(string))
	d.Set("disk", data["disk"].(string))
	d.Set("price", data["price"].(string))
	d.Set("lock_node", data["is_locked"].(bool))
	d.Set("public_ip_address", data["public_ip_address"].(string))
	d.Set("private_ip_address", data["private_ip_address"].(string))
	d.Set("is_bitninja_license_active", data["is_bitninja_license_active"].(bool))
	d.Set("ssh_keys", copy_ssh_keys)
	d.Set("vm_id", int(data["vm_id"].(float64)))

	log.Printf("[info] node Resource read | after setting data")
	if d.Get("status").(string) == "Running" || d.Get("status").(string) == "Creating" {
		d.Set("power_status", "power_on")
	}
	if d.Get("status").(string) == "Powered off" {
		d.Set("power_status", "power_off")
	}
	response, err := apiClient.GetSecurityGroupList(d.Get("project_id").(string), d.Get("location").(string))
	if err != nil {
		log.Printf("[ERROR] Error getting Security Group List inside Node Read. Error : %s", err)
		return diag.Errorf("please confirm the project_id or location that you defined.")
	}
	defaultSG := getDefaultSG(response)
	d.Set("default_sg", defaultSG)

	return diags

}

func resourceUpdateNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)

	nodeId := d.Id()
	project_id := d.Get("project_id").(string)
	location := d.Get("location").(string)
	_, err := apiClient.GetNode(nodeId, project_id)
	if err != nil {

		return diag.Errorf("error finding Item with ID %s", nodeId)

	}

	if d.HasChange("name") {
		log.Printf("[INFO] ndoeId = %v, name = %s ", d.Id(), d.Get("name").(string))
		_, err := apiClient.UpdateNode(nodeId, "rename", d.Get("name").(string), project_id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("power_status") {
		nodestatus := d.Get("status").(string)
		if nodestatus == "Creating" || nodestatus == "Reinstalling" {
			prevBlockIDArray, _ := d.GetChange("block_storage_ids")
			d.Set("block_storage_ids", prevBlockIDArray)
			return diag.Errorf("Node is in %s state", d.Get("status").(string))
		}
		if d.Get("lock_node").(bool) == true {
			return diag.Errorf("cannot change the power status as the node is locked")
		}
		log.Printf("[INFO] %s ", d.Get("power_status").(string))
		apiClient.UpdateNode(nodeId, d.Get("power_status").(string), d.Get("name").(string), project_id)
	}

	if d.HasChange("lock_node") {
		if d.Get("status").(string) == "Creating" || d.Get("status").(string) == "Reinstalling" {
			return diag.Errorf("Cannot update as the node is in %s state", d.Get("status").(string))
		}
		if d.Get("lock_node").(bool) == true {
			_, err := apiClient.UpdateNode(nodeId, "lock_vm", d.Get("name").(string), project_id)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if d.Get("lock_node").(bool) == false {
			_, err := apiClient.UpdateNode(nodeId, "unlock_vm", d.Get("name").(string), project_id)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("reboot_node") {

		if d.Get("reboot_node").(bool) == true {
			d.Set("reboot_node", false)
			if d.Get("status").(string) == "Creating" || d.Get("status").(string) == "Reinstalling" {
				return diag.Errorf("Cannot update as the node is in %s state", d.Get("status").(string))
			}
			if d.Get("status").(string) == "Powered off" {
				return diag.Errorf("cannot reboot as the node is powered off")
			}
			_, err := apiClient.UpdateNode(nodeId, "reboot", d.Get("name").(string), project_id)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("reinstall_node") {
		if d.Get("status").(string) == "Creating" {
			return diag.Errorf("Node is in creating state")
		}
		if d.Get("status").(string) == "Reinstalling" {
			return diag.Errorf("Node already in Reinstalling state")
		}
		if d.Get("reinstall_node").(bool) == true {
			if d.Get("status").(string) == "Powered off" {
				d.Set("reinstall_node", false)
				return diag.Errorf("cannot reinstall as the node is powered off")
			}
			if d.Get("status").(string) == "Reinstalling" {
				d.Set("reinstall_node", false)
				return diag.Errorf("Node already in Reinstalling state")
			}
			_, err := apiClient.UpdateNode(nodeId, "reinstall", d.Get("name").(string), project_id)
			d.Set("reinstall_node", false)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("save_image") {
		if d.Get("save_image") == true {
			d.Set("save_image", false)
			if d.Get("save_image_name").(string) == "" {
				return diag.Errorf("save_image_name empty")
			}

			_, err := apiClient.UpdateNode(nodeId, "save_images", d.Get("save_image_name").(string), project_id)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("security_group_ids") {
		oldSGData, newSGData := d.GetChange("security_group_ids")
		if d.Get("status").(string) != "Running" {
			d.Set("security_group_ids", oldSGData)
			return diag.Errorf("Can only update security groups once the node comes to the running state")
		}
		vm_id := d.Get("vm_id").(int)
		security_groups_list := d.Get("security_group_ids").([]interface{})

		if len(security_groups_list) <= 0 {
			d.Set("security_group_ids", oldSGData)
			return diag.Errorf("Atleast one security groups must be attached to a node!")
		}
		oldSGList := oldSGData.([]interface{})
		newSGList := newSGData.([]interface{})
		sgMap := make(map[int]int)
		for _, sgID := range newSGList {
			sgMap[sgID.(int)] = 1
		}
		for _, sgID := range oldSGList {
			if count, ok := sgMap[sgID.(int)]; ok {
				sgMap[sgID.(int)] = count - 1
			} else {
				sgMap[sgID.(int)] = -1
			}
		}
		var toBeAttached []int
		for key, value := range sgMap {
			if value == -1 {
				log.Printf("----------HAVE TO DETACH THE SECURITY GROUP WITH ID %+v ------------------", key)
				payload := models.UpdateSecurityGroups{
					SecurityGroupList: []int{key},
				}

				response, err := apiClient.DetachSecurityGroup(&payload, vm_id, d.Get("project_id").(string), d.Get("location").(string))
				if err != nil {
					return diag.FromErr(err)
				}
				if _, codeOK := response["code"]; !codeOK {
					return diag.Errorf(response["message"].(string))
				}
				continue
			}
			if value >= 1 {
				toBeAttached = append(toBeAttached, key)
			}
		}
		if len(toBeAttached) >= 1 {
			payload := models.UpdateSecurityGroups{
				SecurityGroupList: toBeAttached,
			}
			response, err := apiClient.AttachSecurityGroup(&payload, vm_id, d.Get("project_id").(string), d.Get("location").(string))
			if err != nil {
				return diag.FromErr(err)
			}
			if _, codeOK := response["code"]; !codeOK {
				return diag.Errorf(response["message"].(string))
			}
		}
	}

	if d.HasChange("label") {
		log.Printf("[INFO] nodeId = %v changed label = %s ", d.Id(), d.Get("label").(string))
		_, err = apiClient.UpdateNode(nodeId, "label_rename", d.Get("label").(string), project_id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("ssh_keys") {
		prevSshKeys, currSshKeys := d.GetChange("ssh_keys")

		log.Printf("[INFO] nodeId = %v changed ssh_keys = %s ", d.Id(), d.Get("ssh_keys"))
		log.Printf("[INFO] type of ssh_keys data = %T", d.Get("ssh_keys"))

		new_SSH_keys, Err := convertLabelToSshKey(m, d.Get("ssh_keys").([]interface{}), project_id)
		if Err != nil {
			d.Set("ssh_keys", prevSshKeys)
			return Err
		}
		d.Set("ssh_keys", new_SSH_keys)
		_, err = apiClient.UpdateNodeSSH(nodeId, "add_ssh_keys", d.Get("ssh_keys").([]interface{}), project_id, d.Get("location").(string))
		d.Set("ssh_keys", currSshKeys)
		if err != nil {
			d.Set("ssh_keys", prevSshKeys)
			return diag.FromErr(err)
		}

	}
	if d.HasChange("location") {
		prevLocation, currLocation := d.GetChange("location")
		log.Printf("[INFO] prevLocation %s, currLocation %s", prevLocation.(string), currLocation.(string))
		d.Set("location", prevLocation)
		return diag.Errorf("location cannot be updated once you create the node.")
	}
	if d.HasChange("image") {
		prevImage, currImage := d.GetChange("image")
		log.Printf("[INFO] prevImage %s, currImage %s", prevImage.(string), currImage.(string))
		d.Set("image", prevImage.(string))
		return diag.Errorf("Image cannot be updated once you create the node.")
	}
	if d.HasChange("plan") {
		prevPlan, currPlan := d.GetChange("plan")

		if d.HasChange("power_status") {
			waitForPoweringOffOn(m, nodeId, project_id)
		}

		log.Printf("[INFO] prevPlan %s, currPlan %s", prevPlan.(string), currPlan.(string))

		if d.Get("status").(string) != "Powered off" {
			d.Set("plan", prevPlan)
			return diag.Errorf("cannot Upgrade as the node is not powered off")
		}
		_, err = apiClient.UpgradeNodePlan(nodeId, d.Get("plan").(string), d.Get("image").(string), project_id, location)

		if err != nil {
			d.Set("plan", prevPlan)
			return diag.FromErr(err)
		}
	}

	if d.HasChange("block_storage_ids") {

		log.Printf("[INFO] Power_status changeing is = %v", d.HasChange("power_status"))
		if d.HasChange("power_status") {
			err := waitForPoweringOffOn(m, nodeId, project_id)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		prevBlockIDArray, currBlockIDArray := d.GetChange("block_storage_ids")

		if d.Get("plan").(string)[0:2] == "C2" {
			d.Set("block_storage_ids", prevBlockIDArray)
			return diag.Errorf("Block storage can not be attached to C2 plan")
		}

		detachingIDs := UniqueArrayElements(prevBlockIDArray.([]interface{}), currBlockIDArray.([]interface{}))
		attachingIDs := UniqueArrayElements(currBlockIDArray.([]interface{}), prevBlockIDArray.([]interface{}))
		CommonIDs := detachingIDs
		log.Printf("[INFO] detachingIDs %+v, attachingIDs %+v, CommonIDs %+v", detachingIDs, attachingIDs, CommonIDs)
		log.Printf("[INFO] prevIDArray %v, currIDArray %v", prevBlockIDArray, currBlockIDArray)

		blockStorage := models.BlockStorageAttach{
			VM_ID: d.Get("vm_id").(int),
		}
		project_id_int, Err := strconv.Atoi(project_id)
		if Err != nil {
			d.Set("block_storage_ids", prevBlockIDArray)
			return diag.FromErr(Err)
		}

		for i, detachingID := range detachingIDs {

			blockStorageID := detachingID.(string)
			_, err := apiClient.AttachOrDetachBlockStorage(&blockStorage, "detach", blockStorageID, project_id_int, location)
			if err != nil {
				d.Set("block_storage_ids", CommonIDs)
				return diag.FromErr(err)
			}
			CommonIDs = removeArrayElement(CommonIDs, detachingID)
			// Wait for some time before detaching the next block storage
			// waitForPoweringOffOn(m, nodeId, project_id)
			if i == len(detachingIDs)-1 {
				break
			}
			waitForDesiredState(apiClient, nodeId, project_id, location)
		}
		for i, attachingID := range attachingIDs {
			blockStorageID := attachingID.(string)
			_, err := apiClient.AttachOrDetachBlockStorage(&blockStorage, "attach", blockStorageID, project_id_int, location)
			if err != nil {
				d.Set("block_storage_ids", CommonIDs)
				log.Printf("[ERROR] Error attaching block storage CommonIDs = %+v", CommonIDs)
				return diag.FromErr(err)
			}
			CommonIDs = append(CommonIDs, attachingID)
			// Wait for some time before attaching the next block storage
			// waitForPoweringOffOn(m, nodeId, project_id)
			if i == len(attachingIDs)-1 {
				break
			}
			waitForDesiredState(apiClient, nodeId, project_id, location)
		}
	}

	return resourceReadNode(ctx, d, m)

}

func resourceDeleteNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	nodeId := d.Id()
	project_id := d.Get("project_id").(string)
	node_status := d.Get("status").(string)
	if node_status == "Saving" || node_status == "Creating" {
		return diag.Errorf("Node in %s state", node_status)
	}
	err := apiClient.DeleteNode(nodeId, project_id, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceExistsNode(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	nodeId := d.Id()
	project_id := d.Get("project_id").(string)
	_, err := apiClient.GetNode(nodeId, project_id)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func convertStringArrayToIntArray(strs []interface{}) ([]int, error) {
	var res []int
	for _, v := range strs {
		i, err := convertStringToInt(v.(string))
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}
func convertStringToInt(str string) (int, error) {
	if str == "" {
		return 0, nil
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func waitForPoweringOffOn(m interface{}, nodeId string, project_id string) error {
	apiClient := m.(*client.Client)

	for {
		// Wait for some time before checking the status again (is Node powered on or off?)
		time.Sleep(constants.WAIT_TIMEOUT * time.Second)

		nodeInfo, err := apiClient.GetNode(nodeId, project_id)
		if err != nil {
			log.Printf("[ERROR] Error getting Node Info inside Plan Upgrade. Error : %s", err)
			return err
		}
		data := nodeInfo["data"].(map[string]interface{})
		log.Printf("[INFO] Node Status : %s", data["status"])
		if data["status"] == "Running" || data["status"] == "Powered off" {
			break
		}
		log.Printf("[INFO] Waiting for Node to power off/on before upgrading the plan")
	}
	return nil
}

func getDefaultSG(response map[string]interface{}) int {
	var res int
	data := response["data"].([]interface{})
	for _, sg := range data {
		sgMap := sg.(map[string]interface{})
		defaultStatus := sgMap["is_default"].(bool)
		if defaultStatus {
			res = int(sgMap["id"].(float64))
			break
		}

	}
	log.Printf("------------Default security group is: %+v -------------", res)
	return res
}

func UniqueArrayElements(arr1 []interface{}, arr2 []interface{}) []interface{} {
	var res []interface{}
	for _, v := range arr1 {
		if !isContains(arr2, v) {
			res = append(res, v)
		}
	}
	return res
}

func isContains(arr []interface{}, val interface{}) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func CommonArrayElements(arr1 []interface{}, arr2 []interface{}) []interface{} {
	var res []interface{}
	for _, v := range arr1 {
		if isContains(arr2, v) {
			res = append(res, v)
		}
	}
	return res
}

func removeArrayElement(arr []interface{}, val interface{}) []interface{} {
	var res []interface{}
	for _, v := range arr {
		if v != val {
			res = append(res, v)
		}
	}
	return res
}

func waitForDesiredState(apiClient *client.Client, nodeId string, project_id string, location string) diag.Diagnostics {
	for {
		// Wait for some time before checking the status again (is Volume Detached?)
		time.Sleep(constants.WAIT_TIMEOUT * time.Second)

		response, err := apiClient.CheckNodeLCMState(nodeId, project_id, location)
		if err != nil {
			log.Printf("[ERROR] Error getting lcm_state %s", err)
			return diag.FromErr(err)
		}
		data := response["data"].(map[string]interface{})
		log.Printf("[INFO] waitForDesiredState data : %+v", data)
		if !(data["lcm_state"].(string) == "Hotplug" || data["lcm_state"].(string) == "HOTPLUG_PROLOG_POWEROFF" || data["lcm_state"].(string) == "HOTPLUG_EPILOG_POWEROFF") {
			break
		}
	}
	return nil
}

func ValidatePlanName(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected plan to be string"))
		return warns, errs
	}
	if value == "" {
		errs = append(errs, fmt.Errorf("plan name cannot be empty"))
		return warns, errs
	}

	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("plan cannot contain whitespace. got %s", value))
		return warns, errs
	}
	return warns, errs
}

func ValidateBlank(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected %s to be string", k))
		return warns, errs
	}
	stripped := strings.TrimSpace(value)
	if stripped == "" {
		errs = append(errs, fmt.Errorf("%s cannot be blank", k))
		return warns, errs
	}
	return warns, errs
}
