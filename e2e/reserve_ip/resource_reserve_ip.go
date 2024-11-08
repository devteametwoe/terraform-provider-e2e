package reserve_ip

import (
	"context"
	// "fmt"
	"log"
	"math"
	"strings"

	// "regexp"
	"strconv"

	//"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceReserveIP() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "specify the project id in which the reserve ip is to be created",
				ForceNew:    true,
			},
			"location": {
				Type:     schema.TypeString,
				Default:  "Delhi",
				Optional: true,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ip address of the reserve ip",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "status of the reserve ip",
			},
			"bought_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "time at which the reserve ip is bought",
			},
			"vm_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "specify the vm id",
			},
			"vm_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vm name",
			},
			"reserve_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "reserve id",
			},
			"appliance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of appliance",
			},
			"reserved_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of reserve ip",
			},
			"project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "project name",
			},
		},

		CreateContext: resourceCreateReserveIP,
		ReadContext:   resourceReadReserveIP,
		DeleteContext: resourceDeleteReserveIP,
		UpdateContext: resourceUpdateReserveIP,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func convertToString(data map[string]interface{}, key string) string {
	if data[key] != nil {
		return data[key].(string)
	}
	return ""
}

func resourceCreateReserveIP(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] NODE CREATE STARTS ")

	res, err := apiClient.NewReservedIp(d.Get("project_id").(string), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] ReservedIp CREATE | RESPONSE BODY | %+v", res)
	if _, codeok := res["code"]; !codeok {
		return diag.Errorf(res["message"].(string))
	}

	if res["is_limit_available"] == false {
		return diag.Errorf(res["message"].(string))
	}

	data := res["data"].(map[string]interface{})
	log.Printf("[INFO] ReserveIP creation | before setting fields %v", data)
	reserveId := data["reserve_id"].(float64)
	reserveId = math.Round(reserveId)
	d.SetId(strconv.Itoa(int(math.Round(reserveId))))
	d.Set("ip_address", convertToString(data, "ip_address"))
	d.Set("status", convertToString(data, "status"))
	d.Set("bought_at", convertToString(data, "bought_at"))
	d.Set("vm_id", convertToString(data, "vm_id"))
	d.Set("vm_name", convertToString(data, "vm_name"))
	d.Set("reserve_id", strconv.FormatFloat(reserveId, 'f', -1, 64))
	d.Set("appliance_type", convertToString(data, "appliance_type"))
	d.Set("reserved_type", convertToString(data, "reserved_type"))
	d.Set("project_name", convertToString(data, "project_name"))
	return diags
}

func resourceReadReserveIP(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	reserveId := d.Get("ip_address").(string)
	project_id := d.Get("project_id").(string)
	res, err := apiClient.GetReservedIp(reserveId, project_id, d.Get("location").(string))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.Errorf("error finding Item with ID %s", reserveId)

		}
	}

	log.Printf("[INFO] ReserveIP READ | RESPONSE BODY | %+v %T", res, res)
	codeok := (res.Code == 200)
	if !codeok {
		return diag.Errorf(res.Message)
	}
	if true || len(res.Data) == 1 {
		data := res.Data[0]
		log.Printf("[INFO] ReserveIP READ | BEFORE SETTING DATA %+v, %v, %T =======================", data, data.Status, data.Status)
		d.Set("ip_address", data.IPAddress)
		d.Set("status", data.Status)
		d.Set("bought_at", data.BoughtAt)
		d.Set("vm_id", data.VMID)
		d.Set("vm_name", data.VMName)
		d.Set("appliance_type", data.ApplianceType)
		d.Set("reserved_type", data.ReservedType)
		d.Set("reserve_id", strconv.FormatFloat(data.ReserveID, 'f', -1, 64))
		d.Set("project_name", data.ProjectName)
	}
	return diags
}

func resourceUpdateReserveIP(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// apiClient := m.(*client.Client)
	// reserveId := d.Id()
	// project_id := d.Get("project_id").(string)
	// _, err := apiClient.GetReservedIp(reserveId, project_id, "Delhi")
	// if err != nil || err == nil {
	// 	return diag.Errorf("*******----Cannot Update reserve ip----*******")
	// }
	// return diag.Errorf("cannot update reserve ip")
	return diags
}

func resourceDeleteReserveIP(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	ip_address := d.Get("ip_address").(string)
	project_id := d.Get("project_id").(string)

	err := apiClient.DeleteReserveIP(ip_address, project_id, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
