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

	//"time"
	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the Plan",
			},
			"backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Tells you the state of your backups",
				Default:     false,
			},

			"image": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the image you have selected format :- ( os-version )",
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
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "region",
				Default:     "ncr",
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
			"security_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specify the security group. Checkout security_groups datasource listing security groups",
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
			},
			"vm_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The id of the VM.",
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
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
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

	log.Printf("[INFO] NODE CREATE STARTS ")
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
		Region:            d.Get("region").(string),
		Reserve_ip:        d.Get("reserve_ip").(string),
		Vpc_id:            d.Get("vpc_id").(string),
		Security_group_id: d.Get("security_group_id").(int),
		SSH_keys:          d.Get("ssh_keys").([]interface{}),
		Start_scripts:     d.Get("start_scripts").([]interface{}),
	}

	if node.Vpc_id != "" {
		vpc_details, err := apiClient.GetVpc(node.Vpc_id, d.Get("project_id").(string), d.Get("region").(string))
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
	if d.Get("status").(string) == "Running" {
		d.Set("power_status", "power_on")
	}
	if d.Get("status").(string) == "Powered off" {
		d.Set("power_status", "power_off")
	}

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
