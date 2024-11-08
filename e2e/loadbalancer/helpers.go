package loadbalancer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetLbPort(mode string) string {
	if mode == "HTTP" {
		return "80"
	}
	return "443"
}

func ExpandBackends(config []interface{}, apiClient *client.Client, project_id string) ([]models.Backend, error) {
	backends := make([]models.Backend, 0, len(config))

	for _, backend := range config {
		detail := backend.(map[string]interface{})

		servers, err := ExpandServers(detail["servers"].(interface{}), apiClient, project_id)
		if err != nil {
			return nil, err
		}
		r := models.Backend{
			Balance:        detail["balance"].(string),
			CheckboxEnable: detail["checkbox_enable"].(bool),
			DomainName:     detail["domain_name"].(string),
			CheckUrl:       detail["check_url"].(string),
			Servers:        servers,
			HttpCheck:      detail["http_check"].(bool),
			Name:           detail["name"].(string),
			ScalerId:       detail["scaler_id"].(string),
			ScalerPort:     detail["scaler_port"].(string),
		}
		backends = append(backends, r)
	}
	return backends, nil
}

func ExpandServers(server_details interface{}, apiClient *client.Client, project_id string) ([]models.Server, error) {
	var servers []models.Server

	for _, server := range server_details.([]interface{}) {
		server_detail := server.(map[string]interface{})
		node, err := apiClient.GetNode(server_detail["id"].(string), project_id)
		if err != nil {
			return nil, err
		}
		data := node["data"].(map[string]interface{})
		status := data["status"].(string)
		if status != "Running" {
			return nil, fmt.Errorf("Node with id %s is not in running state", server_detail["id"].(string))
		}
		r := models.Server{
			BackendName: data["name"].(string),
			BackendIp:   data["private_ip_address"].(string),
			BackendPort: server_detail["port"].(string),
		}

		servers = append(servers, r)
	}
	if len(servers) == 0 {
		return make([]models.Server, 0), nil
	}
	return servers, nil
}

func ExpandAclList(config []interface{}) ([]models.AclListInfo, error) {
	aclList := make([]models.AclListInfo, 0, len(config))

	for _, acl := range config {
		aclRule := acl.(map[string]interface{})

		r := models.AclListInfo{
			AclName:         aclRule["acl_name"].(string),
			AclCondition:    aclRule["acl_condition"].(string),
			AclMatchingPath: aclRule["acl_matching_path"].(string),
		}

		aclList = append(aclList, r)
	}
	return aclList, nil
}

func ExpandAclMap(config []interface{}) ([]models.AclMapInfo, error) {
	aclMap := make([]models.AclMapInfo, 0, len(config))

	for _, backendlist := range config {
		aclMapData := backendlist.(map[string]interface{})

		r := models.AclMapInfo{
			AclName:           aclMapData["acl_name"].(string),
			AclConditionState: true,
			AclBackend:        aclMapData["acl_backend"].(string),
		}

		aclMap = append(aclMap, r)
	}
	return aclMap, nil
}

func ExpandVpcList(d *schema.ResourceData, vpc_list []interface{}, apiClient *client.Client) ([]models.VpcDetail, error) {
	var vpc_details []models.VpcDetail

	for _, id := range vpc_list {
		vpc_detail, err := apiClient.GetVpc(strconv.Itoa(id.(int)), d.Get("project_id").(string), d.Get("location").(string))
		if err != nil {
			return nil, err
		}
		data := vpc_detail.Data
		if data.State != "Active" {
			return nil, fmt.Errorf("Can not attach vpc currently, vpc is in %s state", data.State)
		}
		r := models.VpcDetail{
			Network_id: data.Network_id,
			VpcName:    data.Name,
			Ipv4_cidr:  data.Ipv4_cidr,
		}

		vpc_details = append(vpc_details, r)
	}
	return vpc_details, nil
}

func ExpandEnableEosLogger(config []interface{}) (models.EosDetail, error) {
	eosDetail := models.EosDetail{}

	for _, eosBucketInfo := range config {
		detail := eosBucketInfo.(map[string]interface{})
		eosDetail.ApplianceId = detail["appliance_id"].(int)
		eosDetail.AccessKey = detail["access_key"].(string)
		eosDetail.Secretkey = detail["secret_key"].(string)
		eosDetail.Bucket = detail["bucket"].(string)
	}
	return eosDetail, nil
}

func ExpandTcpBackend(config []interface{}, apiClient *client.Client, project_id string) ([]models.TcpBackendDetail, error) {
	tcpBackends := make([]models.TcpBackendDetail, 0, len(config))

	for _, tcpBackend := range config {
		detail := tcpBackend.(map[string]interface{})

		servers, err := ExpandServers(detail["servers"].(interface{}), apiClient, project_id)
		if err != nil {
			return nil, err
		}
		r := models.TcpBackendDetail{
			BackendName: detail["backend_name"].(string),
			Port:        detail["port"].(string),
			Balance:     detail["balance"].(string),
			Servers:     servers,
		}

		tcpBackends = append(tcpBackends, r)
	}
	return tcpBackends, nil
}

func SetLoadBalancerStatus(d *schema.ResourceData, status_detail interface{}) error {
	haproxyStatus := status_detail.(map[string]interface{})
	dataMonitor := haproxyStatus["data_monitor"].(map[string]interface{})
	if haproxyStatus["status"] == "RUNNING" {
		if len(dataMonitor) == 0 {
			d.Set("status", "Backend Status Unavailable")
			return nil
		}
		if dataMonitor["status"].(bool) == false {
			d.Set("status", "Backend Connection Failure")
		} else {
			d.Set("status", "Running")
		}
	} else if haproxyStatus["status"] == "STOP" {
		d.Set("status", "Powered off")
	} else if haproxyStatus["status"] == "Creating" {
		d.Set("status", "Creating")
	} else if haproxyStatus["status"] == "Deploying" {
		d.Set("status", "Deploying")
	} else if haproxyStatus["status"] == "UPDATING" {
		d.Set("status", "Upgrading")
	} else {
		d.Set("status", "Error")
	}
	return nil
}

func CheckStatus(statuslist []string, status string) bool {
	for _, s := range statuslist {
		if strings.EqualFold(s, status) {
			return true
		}
	}
	return false
}
