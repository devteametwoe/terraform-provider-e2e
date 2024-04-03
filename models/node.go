package models

type NodeCreate struct {
	Name                    string        `json:"name"`
	Label                   string        `json:"label"`
	Plan                    string        `json:"plan"`
	Backup                  bool          `json:"backup"`
	Image                   string        `json:"image"`
	Default_public_ip       bool          `json:"default_public_id"`
	Disable_password        bool          `json:"disable_password"`
	Enable_bitninja         bool          `json:"enable_bitninja"`
	Is_ipv6_availed         bool          `json:"is_ipv6_availed"`
	Is_saved_image          bool          `json:"is_saved_image"`
	Region                  string        `json:"region"`
	Reserve_ip              string        `json:"reserve_ip"`
	Vpc_id                  string        `json:"vpc_id"`
	Ngc_container_id        int           `json:"ngc_container_id"`
	Saved_image_template_id int           `json:"saved_image_template_id"`
	Security_group_id       int           `json:"security_group_id"`
	SSH_keys                []interface{} `json:"ssh_keys"`
	Start_scripts           []interface{} `json:"start_scripts"`
	Image_id                int           `json:"image_id"`
}
type NodeAction struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type NodeActionSSH struct {
	Type     string                   `json:"type"`
	SSH_KEYS []map[string]interface{} `json:"ssh_keys"`
}

type NodePlanUpgradeAction struct {
	Plan  string `json:"plan"`
	Image string `json:"image"`
}

type ResponseNodes struct {
	Code    int    `json:"code"`
	Data    []Node `json:"data"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Node struct {
	ID               float64 `json:"id"`
	Name             string  `json:"name"`
	Status           string  `json:"status"`
	PublicIPAddress  string  `json:"public_ip_address"`
	PrivateIPAddress string  `json:"private_ip_address"`
	RescueModeStatus string  `json:"rescue_mode_status"`
	IsLocked         bool    `json:"is_locked"`
}

type UpdateSecurityGroups struct {
	SecurityGroupList []int `json:"security_group_ids"`
}
