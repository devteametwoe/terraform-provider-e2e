package sfs

import (
 "context"
 //"encoding/json"
 // "fmt"
 "log"
 // "math"
 // "regexp"

 // "strconv"
 //"strings"

 "github.com/e2eterraformprovider/terraform-provider-e2e/models"

 // "github.com/hashicorp/terraform-plugin-log"
 // "github.com/hashicorp/terraform-plugin-log/tflog"

 "github.com/e2eterraformprovider/terraform-provider-e2e/client"
 "github.com/hashicorp/terraform-plugin-sdk/v2/diag"

 "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSfs() *schema.Resource {
 return &schema.Resource{
     Schema: map[string]*schema.Schema{
         "region": {
             Type:        schema.TypeString,
             Optional:    true,
             Description: "Region should specified",
         },
         "project_id":{
                Type:        schema.TypeString,
             Required:    true,
             Description: "project_id is mandatory",
         },
         "sfs_list": {
             Type:        schema.TypeList,
             Computed:    true,
             Description: "List of all the SFS of your account . ",
             Elem: &schema.Resource{
                 Schema: map[string]*schema.Schema{
                     "id": {
                         Type:        schema.TypeInt,
                         Computed:    true,
                         Description: "The id of the node",
                     },
                     "name": {
                         Type:     schema.TypeString,
                         Computed: true,
                     },
                     "efs_disk_size": {
                         Type:     schema.TypeString,
                         Computed: true,
                     },
                     "status": {
                         Type:     schema.TypeString,
                         Computed: true,
                     },
                     "private_endpoint": {
                         Type:     schema.TypeString,
                         Computed: true,
                     },
                     "plan_name":{
                         Type:     schema.TypeString,
                         Computed:  true,
                     },
                     "is_backup_enabled":{
                         Type:     schema.TypeBool,
                         Computed:  true,
                     },
                     "iops":{
                         Type:    schema.TypeInt,
                         Computed:  true,
                     },

                 },
             },
         },
     },
     ReadContext: dataSourceReadSfs,
     Importer: &schema.ResourceImporter{
         State: schema.ImportStatePassthrough,
     },
 }
}

func dataSourceReadSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

 var diags diag.Diagnostics
 apiClient := m.(*client.Client)
 log.Printf("[INFO] Inside nodes data source ")
 project_id:=d.Get("project_id").(string)
 Response, err := apiClient.GetSfss(d.Get("region").(string), project_id)
 if err != nil {
     return diag.FromErr(err)
 }
 log.Printf("[INFO] NODES DATA SOURCE | before setting")
 d.Set("sfs_list", flattenSfs(&Response.Data))
 d.SetId("sfs_list")

 return diags
}

func flattenSfs(nodes *[]models.SfssRead) []interface{} {

 if nodes != nil {
     ois := make([]interface{}, len(*nodes), len(*nodes))

     for i, node := range *nodes {
         oi := make(map[string]interface{})
         oi["id"] = node.ID
         oi["name"] = node.Name
         oi["efs_disk_size"]=node.DiskSize
         oi["plan_name"]=node.PlanName
         oi["status"] = node.Status
         oi["private_endpoint"] = node.PrivateIPAddress
         oi["is_backup_enabled"]=node.IsBackup
         oi["iops"]=node.Iops
         ois[i] = oi
     }
     return ois
 }
 return make([]interface{}, 0)
}