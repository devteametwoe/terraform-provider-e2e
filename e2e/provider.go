package e2e

import (
	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/blockstorage"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/image"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/kubernetes"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/loadbalancer"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/objectstore"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/reserve_ip"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/security_group"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sfs"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/ssh_key"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{

			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_API_KEY", ""),
				Description: "valied api key required ",
			},
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_AUTH_TOKEN", ""),
				Description: "Valied authentication Bearer token required",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.e2enetworks.com/myaccount/api/v1/",
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_API_ENDPOINT", "https://api.e2enetworks.com/myaccount/api/v1"),
				Description: "specify the endpoint , default endpoint is https://api.e2enetworks.com/myaccount/api/v1/",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"example_node":         node.ResourceNode(),
			"example_image":        image.ResourceImage(),
			"example_loadbalancer": loadbalancer.ResourceLoadBalancer(),
			"example_vpc":          vpc.ResouceVpc(),
			"example_reserved_ip":  reserve_ip.ResourceReserveIP(),
			"example_blockstorage": blockstorage.ResourceBlockStorage(),
			"example_sfs":          sfs.ResourceSfs(),
			"example_objectstore":  objectstore.ResourceObjectStore(),
			"example_ssh_key":      ssh_key.ResourceSshKey(),
			"example_kubernetes":   kubernetes.ResourceKubernetesService(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"example_node":            node.DataSourceNode(),
			"example_images":          image.DataSourceImages(),
			"example_ssh_keys":        ssh_key.DataSourceSshKeys(),
			"example_ssh_key":         ssh_key.DataSourceSshKey(),
			"example_vpcs":            vpc.DataSourceVpcs(),
			"example_security_groups": security_group.DataSourceSecurityGroups(),
			"example_blockstorage":    blockstorage.DataSourceBlockStorage(),
			"example_nodes":           node.DataSourceNodes(),
			"example_reserve_ips":     reserve_ip.DataSourceReserveIps(),
			"example_sfss":            sfs.DataSourceSfs(),
			"example_objectstores":    objectstore.DataSourceObjectStores(),
			"example_kubernetes":      kubernetes.DataSourceKubernetesService(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	api_key := d.Get("api_key").(string)
	auth_token := d.Get("auth_token").(string)
	api_endpoint := d.Get("api_endpoint").(string)
	return client.NewClient(api_key, auth_token, api_endpoint), nil
}
