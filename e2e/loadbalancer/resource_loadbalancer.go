package loadbalancer

import (
	"context"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var nameRegex string = "^[a-zA-Z0-9-_]{0,50}$"

func ResourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Schema:        ResouceLoadBalancerSchema(),
		CreateContext: resourceCreateLoadBalancer,
		ReadContext:   resourceReadLoadBalancer,
		UpdateContext: resourceUpdateLoadBalancer,
		DeleteContext: resourceDeleteLoadBalancer,
		Exists:        resourceExistsLoadBalancer,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func ResouceLoadBalancerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"plan_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "It is the plan of which load balancer is going to launch",
			ValidateFunc: validation.StringInSlice([]string{
				"E2E-LB-2",
				"E2E-LB-3",
				"E2E-LB-4",
				"E2E-LB-5",
			}, false),
		},
		"lb_name": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "It is the name of load balancer, letter,digit,underscore,hyphen are allowed",
			ValidateFunc: node.ValidateName,
		},
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "This is your project ID in which you want to create the resource.",
			ForceNew:    true,
		},
		"lb_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "External",
			Description: "It is used to define internal or extenal load balancer",
		},
		"lb_mode": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "In which mode load balancer is going to launch http, https, both",
			ValidateFunc: validation.StringInSlice([]string{
				"HTTP",
				"HTTPS",
				"Both",
			}, false),
		},
		"node_list_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "S",
			Description: "It is used to find out either node is static or dynamic autoscaling",
			ValidateFunc: validation.StringInSlice([]string{
				"S",
				"D",
			}, false),
		},
		"checkbox_enable": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "This checkbox is", // need description for this checkbox
		},
		"lb_reserve_ip": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "This field is for any reserve IP which is going to attach on load balancer",
		},
		"ssl_certificate_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "This field is used to set ssl sertificate if lb mode is https or both",
		},
		"ssl_context": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "This field is used to set ssl context",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"redirect_to_https": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "If Load balancer is set to both http and https then this option need to select",
					},
				},
			},
		},
		"enable_bitninja": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Modular security tool used to enable load balancer from wide range of cyber attacks",
		},
		"backends": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "This will contain the backend details which will be attached to load balancer",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "This will be the name of your backend.",
					},
					"scaler_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "",
						Description: "Need scalar ID if you want to attach autoscaling",
					},
					"scaler_port": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "",
						Description: "Need scalar port if you want to attach autoscaling",
					},
					"balance": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "This will contain the type of algorithim used while load balancing",
						ValidateFunc: validation.StringInSlice([]string{
							"source",
							"roundrobin",
							"leastconn",
						}, false),
					},
					"checkbox_enable": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "This checkbox is to enable healthcheck",
					},
					"domain_name": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "localhost",
						Description: "domain name for healthcheck",
					},
					"check_url": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "/",
						Description: "endpoint of healthckeck to ping",
					},
					"servers": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "description of servers that are going to attach on backend",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Node id which you want to attach",
								},
								"port": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Port Number of the node",
								},
							},
						},
					},
					"http_check": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Check if http health check in enable",
					},
				},
			},
		},
		"acl_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "This will give the acl rule which you want to apply",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"acl_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of your ACL rule",
					},
					"acl_condition": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Condition in which ACL rule will match",
					},
					"acl_matching_path": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "path in which this rule will work",
					},
				},
			},
		},

		"acl_map": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "This will give you how you want to route request according to acl rule",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"acl_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of your ACL rule",
					},
					"acl_condition_state": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     true,
						Description: "status of acl condition state",
					},
					"acl_backend": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of your backend server",
					},
				},
			},
		},
		"vpc_list": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Optional:    true,
			Description: "List of vpc Id which you want to attach",
		},
		"enable_eos_logger": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "If you want to get the logs of loadbalancer. Please connect eos bucket",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"appliance_id": {
						Type:        schema.TypeInt,
						Optional:    true,
						Default:     0,
						Description: "ID of the appliance",
					},
					"access_key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Access key of your object storage bucket",
					},
					"secret_key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Secret key of your object storage bucket",
					},
					"bucket": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Bucket name of your object storage bucket",
					},
				},
			},
		},
		"tcp_backend": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Need Information of TCP backend If user want to attach",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"backend_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Your TCP backend name",
						ValidateFunc: node.ValidateName,
					},
					"port": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Port number for your TCP backend. 8080, 10050, 9101,80 or 443 port not allowed",
						ValidateFunc: validation.StringNotInSlice([]string{
							"8080",
							"10050",
							"9101",
							"80",
							"443",
						}, false),
					},
					"balance": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "This will contain the type of algorithim used while load balancing",
						ValidateFunc: validation.StringInSlice([]string{
							"source",
							"roundrobin",
							"leastconn",
						}, false),
					},
					"servers": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "description of servers that are going to attach on backend",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Type:         schema.TypeString,
									Required:     true,
									Description:  "Node id which you want to attach",
									ValidateFunc: node.ValidateName,
								},
								"port": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Port Number of the node",
								},
							},
						},
					},
				},
			},
		},
		"is_ipv6_attached": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "This is used to attach IPV6 on your load balancer",
		},
		"default_backend": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"power_status": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "power_on",
			Description: "power_on to start the load balancer and power_off to power off the load balancer",
		},
		"public_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Public IP of load balancer",
		},
		"private_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Private IP of load balancer",
		},
		"ram": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "This is the ram allotted to your loadbalancer",
		},
		"disk": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "This is the disk storage allotted to your loadbalancer",
		},
		"vcpu": {
			Type:        schema.TypeFloat,
			Computed:    true,
			Description: "This is the vcpu allotted to your loadbalancer",
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Delhi",
			Description: "This is the region of your loadbalancer",
			ForceNew:    true,
		},
		"host_target_ipv6": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "This is the ipv6 allotted to your loadbalancer",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "This is the status of your loadbalancer, only to get the status from my account.",
		},
	}
}

func CreateLoadBalancerObject(apiClient *client.Client, d *schema.ResourceData) (*models.LoadBalancerCreate, diag.Diagnostics) {
	log.Printf("[INFO] LOAD BALANCER OBJECT CREATION STARTS")

	loadBalancerObj := models.LoadBalancerCreate{
		PlanName:         d.Get("plan_name").(string),
		LbName:           d.Get("lb_name").(string),
		LbType:           d.Get("lb_type").(string),
		LbMode:           d.Get("lb_mode").(string),
		LbPort:           GetLbPort(d.Get("lb_mode").(string)),
		NodeListType:     d.Get("node_list_type").(string),
		CheckBoxEnable:   d.Get("checkbox_enable").(string),
		LbReserveIp:      d.Get("lb_reserve_ip").(string),
		SslCertificateId: d.Get("ssl_certificate_id").(string),
		EnableBitninja:   d.Get("enable_bitninja").(bool),
		IsIpv6Attached:   d.Get("is_ipv6_attached").(bool),
		DefaultBackend:   d.Get("default_backend").(string),
	}
	enableEosLogger, ok := d.GetOk("enable_eos_logger")
	if ok {
		eosDetail, err := ExpandEnableEosLogger(enableEosLogger.(*schema.Set).List())
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.EnableEosLogger = eosDetail
	}
	aclList, ok := d.GetOk("acl_list")
	if ok {
		aclListDetail, err := ExpandAclList(aclList.([]interface{}))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.AclList = aclListDetail
	} else {
		loadBalancerObj.AclList = make([]models.AclListInfo, 0)
	}
	aclMap, ok := d.GetOk("acl_map")
	if ok {
		aclMapDetail, err := ExpandAclMap(aclMap.([]interface{}))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.AclMap = aclMapDetail
	} else {
		loadBalancerObj.AclMap = make([]models.AclMapInfo, 0)
	}
	tcpBackend, ok := d.GetOk("tcp_backend")
	if ok {
		tcpBackendDetail, err := ExpandTcpBackend(tcpBackend.([]interface{}), apiClient, d.Get("project_id").(string))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.TcpBackend = tcpBackendDetail
	} else {
		loadBalancerObj.TcpBackend = make([]models.TcpBackendDetail, 0)
	}

	backends, ok := d.GetOk("backends")
	if ok {
		backendDetail, err := ExpandBackends(backends.([]interface{}), apiClient, d.Get("project_id").(string))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.Backends = backendDetail
	} else {
		loadBalancerObj.Backends = make([]models.Backend, 0)
	}

	vpcList, ok := d.GetOk("vpc_list")
	if ok {
		vpcListDetail, err := ExpandVpcList(d, vpcList.(*schema.Set).List(), apiClient)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.VpcList = vpcListDetail
	} else {
		loadBalancerObj.VpcList = make([]models.VpcDetail, 0)
	}

	sslContext, ok := d.GetOk("ssl_context")
	if ok {
		sslContextList := sslContext.([]interface{})
		detail := sslContextList[0].(map[string]interface{})
		loadBalancerObj.SslContext = detail
	} else {
		loadBalancerObj.SslContext = map[string]interface{}{"redirect_to_https": false}
	}
	return &loadBalancerObj, nil
}
func resourceCreateLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	loadBalancerObj, diags := CreateLoadBalancerObject(apiClient, d)
	if diags != nil {
		return diags
	}
	response, err := apiClient.NewLoadBalancer(loadBalancerObj, d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] LOAD BALANCER CREATE | RESPONSE BODY | %+v", response)

	if _, codeok := response["code"]; !codeok {
		return diag.Errorf(response["message"].(string))
	}

	data := response["data"].(map[string]interface{})
	if data["is_credit_sufficient"] == false {
		return diag.Errorf("Credit is not sufficient")
	}
	log.Printf("[INFO] load balancer creation | before setting fields")

	lbId := data["id"].(float64)
	lbId = math.Round(lbId)
	d.SetId(strconv.Itoa(int(math.Round(lbId))))
	d.Set("public_ip", data["IP"].(string))
	return diags
}

func resourceReadLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("=============INSIDE RESOURCE READ LOAD BALANCER==========================")
	lbId := d.Id()
	location := d.Get("location").(string)
	lb, err := apiClient.GetLoadBalancerInfo(lbId, location, d.Get("project_id").(string))
	log.Println("===========GET_LOAD_BALANCER_RESPONSE==========", lb)
	if err != nil {
		return diag.Errorf("error finding Item with ID %s", lbId)
	}

	log.Printf("[INFO] LOADBALANCER READ | BEFORE SETTING DATA")
	data := lb["data"].(map[string]interface{})
	node_detail := data["node_detail"].(map[string]interface{})
	appliance_instance := data["appliance_instance"].([]interface{})
	instance := appliance_instance[0].(map[string]interface{})
	lb_context := instance["context"].(map[string]interface{})
	d.Set("private_ip", node_detail["private_ip"].(string))
	d.Set("public_ip", node_detail["public_ip"].(string))
	d.Set("ram", node_detail["ram"].(string))
	d.Set("disk", node_detail["disk"].(string))
	d.Set("vcpu", node_detail["vcpu"].(float64))
	d.Set("lb_name", data["name"].(string))
	d.Set("plan_name", node_detail["plan_name"].(string))
	d.Set("lb_mode", lb_context["lb_mode"].(string))

	if d.Get("is_ipv6_attached").(bool) == true {
		if lb_context["host_target_ipv6"] != nil {
			d.Set("host_target_ipv6", lb_context["host_target_ipv6"].(string))
		} else {
			d.Set("is_ipv6_attached", false)
		}
	}
	err = SetLoadBalancerStatus(d, data["lb_status"])
	if err != nil {
		return diag.Errorf("error while setting lb status with ID %s, error : %s", lbId, err)
	}
	if d.Get("status").(string) == "Powered off" {
		d.Set("power_status", "power_off")
	} else {
		d.Set("power_status", "power_on")
	}
	return diags
}

func resourceUpdateLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	lbId := d.Id()
	location := d.Get("location").(string)
	lb_status := d.Get("status").(string)
	response, err := apiClient.GetLoadBalancerInfo(lbId, location, d.Get("project_id").(string))
	data := response["data"].(map[string]interface{})
	if err != nil {
		return diag.Errorf("error Fetching Load Balancer resource with ID %s", lbId)
	}

	node_detail := data["node_detail"].(map[string]interface{})

	if d.HasChange("power_status") {
		disablePowerStatusList := []string{"Creating", "Deploying", "Upgrading"}

		if CheckStatus(disablePowerStatusList, lb_status) {
			return diag.Errorf("Load Balancer is in %s state, can not change power status.", lb_status)
		}

		payload := map[string]interface{}{"type": d.Get("power_status").(string)}
		err := apiClient.UpdateLoadBalancerAction(payload, lbId, location, d.Get("project_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		return resourceReadLoadBalancer(ctx, d, m)
	}

	if d.HasChange("plan_name") {
		currentPlanName := node_detail["plan_name"].(string)
		newPlanName := d.Get("plan_name").(string)
		if strings.Compare(newPlanName, currentPlanName) == -1 {
			return diag.Errorf("Can not downgrade your plan. Kindly provide the higher plan name")
		}
		payload := map[string]interface{}{
			"type":      "upgrade_plan",
			"name":      data["name"].(string),
			"plan_name": newPlanName,
		}
		err := apiClient.UpdateLoadBalancerAction(payload, lbId, location, d.Get("project_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		return resourceReadLoadBalancer(ctx, d, m)
	}

	if lb_status == "Powered off" {
		return diag.Errorf("Can not Update Load Balancer as it is in %s state", lb_status)
	}

	if d.HasChange("lb_name") {
		payload := map[string]interface{}{
			"type": "rename",
			"name": d.Get("lb_name").(string),
		}
		err := apiClient.UpdateLoadBalancerAction(payload, lbId, location, d.Get("project_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		return resourceReadLoadBalancer(ctx, d, m)
	}

	appliance_instance := data["appliance_instance"].([]interface{})
	instance := appliance_instance[0].(map[string]interface{})
	lb_context := instance["context"].(map[string]interface{})

	if d.HasChange("is_ipv6_attached") {
		ipv6_attach := d.Get("is_ipv6_attached").(bool)
		payload := map[string]interface{}{}
		if ipv6_attach == true {
			payload = map[string]interface{}{"action": "attach"}
		} else {
			payload = map[string]interface{}{
				"action":      "detach",
				"detach_ipv6": lb_context["host_target_ipv6"].(string),
			}
		}
		err := apiClient.IPV6LoadBalancerAction(payload, lbId, location, d.Get("project_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		return resourceReadLoadBalancer(ctx, d, m)
	}

	loadBalancerObj, diags := CreateLoadBalancerObject(apiClient, d)
	if diags != nil {
		return diags
	}
	res, err := apiClient.LoadBalancerBackendUpdate(loadBalancerObj, lbId, location, d.Get("project_id").(string))
	resData := res["data"].(map[string]interface{})
	if resData["is_credit_sufficient"] == false {
		return diag.Errorf("Credit is not sufficient")
	}
	return resourceReadLoadBalancer(ctx, d, m)
}

func resourceDeleteLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	lbId := d.Id()
	lb_status := d.Get("status").(string)
	disableDeleteLbStatusList := []string{"Creating", "Deploying", "Upgrading"}

	if CheckStatus(disableDeleteLbStatusList, lb_status) {
		return diag.Errorf("Load Balancer is in %s state. Currently can not destroy the resource.", lb_status)
	}

	err := apiClient.DeleteLoadBalancer(lbId, d.Get("location").(string), d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceExistsLoadBalancer(d *schema.ResourceData, m interface{}) (bool, error) {
	return true, nil
}

// func ExpandVpcList(d *schema.ResourceData, vpc_list []interface{}, apiClient *client.Client) ([]models.VpcDetail, error) {
// 	var vpc_details []models.VpcDetail

// 	for _, id := range vpc_list {
// 		vpc_detail, err := apiClient.GetVpc(strconv.Itoa(id.(int)), d.Get("project_id").(int), d.Get("location").(string))
// 		if err != nil {
// 			return nil, err
// 		}
// 		data := vpc_detail.Data
// 		if data.State != "Active" {
// 			return nil, fmt.Errorf("Can not attach vpc currently, vpc is in %s state", data.State)
// 		}
// 		r := models.VpcDetail{
// 			Network_id: data.Network_id,
// 			VpcName:    data.Name,
// 			Ipv4_cidr:  data.Ipv4_cidr,
// 		}

// 		vpc_details = append(vpc_details, r)
// 	}
// 	return vpc_details, nil
// }
