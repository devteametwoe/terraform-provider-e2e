package models

type BlockStorageCreate struct {
	Name string  `json:"name"`
	Size float64 `json:"size"`
	IOPS string  `json:"iops"`
}

type BlockStorageUpgrade struct {
	Name  string  `json:"name"`
	Size  float64 `json:"block_storage_size"`
	VM_ID float64 `json:"vm_id"`
}

type BlockStorageResponse struct {
	Code    int                    `json:"code"`
	Data    []BlockStorage         `json:"data"`
	Errors  map[string]interface{} `json:"errors"`
	Message string                 `json:"message"`
}

type ResponseTemplate struct {
	DevPrefix    string `json:"DEV_PREFIX"`
	Driver       string `json:"DRIVER"`
	TotalIOPSSec string `json:"TOTAL_IOPS_SEC"`
}

type BlockStorage struct {
	BlockID   int                    `json:"block_id"`
	Name      string                 `json:"name"`
	Size      int                    `json:"size"`
	Status    string                 `json:"status"`
	Template  ResponseTemplate       `json:"template"`
	VMDetail  map[string]interface{} `json:"vm_detail"`
	CreatedOn string                 `json:"created_on"`
}

type BlockStorageAttach struct {
	VM_ID int `json:"vm_id"`
}
