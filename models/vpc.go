package models

type VpcsResponse struct {
	Code    int           `json:"code"`
	Data    []Vpc         `json:"data"`
	Error   []interface{} `json:"error"`
	Message string        `json:"message"`
}
type Vpc struct {
	Created_at string  `json:"created_at,omitempty"`
	State      string  `json:"state,omitempty"`
	Name       string  `json:"name,omitempty"`
	Ipv4_cidr  string  `json:"ipv4_cidr,omitempty"`
	Network_id float64 `json:"network_id,omitempty"`
	Gateway_ip string  `json:"gateway_ip,omitempty"`
	Pool_size  float64 `json:"pool_size,omitempty"`
	Is_active  bool    `json:"is_active,omitempty"`
}

type VpcResponse struct {
	Code    int           `json:"code"`
	Data    Vpc           `json:"data"`
	Error   []interface{} `json:"error"`
	Message string        `json:"message"`
}

type VpcCreate struct {
	VpcName     string `json:"vpc_name"`
	NetworkSize int    `json:"network_size"`
}
