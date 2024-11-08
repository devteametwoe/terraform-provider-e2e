package models

type ResponseReserveIps struct {
	Code    int         `json:"code"`
	Data    []ReserveIp `json:"data"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
}

type ReserveIp struct {
	IPAddress     string  `json:"ip_address"`
	Status        string  `json:"status"`
	BoughtAt      string  `json:"bought_at"`
	VMID          float64 `json:"vm_id"`
	VMName        string  `json:"vm_name"`
	ReserveID     float64 `json:"reserve_id"`
	ApplianceType string  `json:"appliance_type"`
	ReservedType  string  `json:"reserved_type"`
	ProjectName   string  `json:"project_name"`
}
