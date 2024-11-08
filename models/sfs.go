package models

type SfsCreate struct{
	Name                    string        `json:"efs_name"`
	Plan                    string        `json:"efs_plan_name"`
	Vpc_id                  string         `json:"vpc_id"`
	Disk_size               int           `json:"efs_disk_size"`
	Disk_iops               int           `json:"efs_disk_iops"`
}

type SfssRead struct{
	ID               int      `json:"id"`
	Name             string  `json:"name"`
	DiskSize         string   `json:"efs_disk_size"`
	Status           string  `json:"status"`
	PrivateIPAddress string  `json:"private_endpoint"`
	Iops             int     `json:"iops"`
	IsBackup         bool    `json:"is_backup_enabled"`
	PlanName         string   `json:"plan_name"`
}

type ResponseSfss struct {
	Code    int    `json:"code"`
	Data    []SfssRead `json:"data"`
	Error   string `json:"error"`
	Message string `json:"message"`
}