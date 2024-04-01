package node

import (
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func convertLabelToSshKey(m interface{}, ssh_keys []interface{}, project_id string) ([]interface{}, diag.Diagnostics) {

	apiClient := m.(*client.Client)

	log.Printf("[INFO] Helper Function ssh_keys = %+v", ssh_keys)
	if ssh_keys != nil || len(ssh_keys) > 0 {
		var new_SSH_keys []interface{}
		for _, v := range ssh_keys {
			res, err := apiClient.GetSshKey(v.(string), project_id)
			log.Printf("[INFO] Helper Function res = %+v", res)
			if err != nil {
				return nil, diag.FromErr(err)
			}
			if code, codeok := res["code"].(float64); !codeok || int(code) < 200 || int(code) >= 300 {
				log.Printf("code and codeok, %v, %v", code, codeok)
				return nil, diag.Errorf("%+v", res["errors"])
			}
			data := res["data"].(map[string]interface{})
			new_SSH_keys = append(new_SSH_keys, data["ssh_key"].(string))
		}
		return new_SSH_keys, nil
	}
	return nil, nil
}
