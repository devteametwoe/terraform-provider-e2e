package models

type ImageListResponse struct {
	Code    int     `json:"code"`
	Data    []Image `json:"data"`
	Error   string  `json:"error"`
	Message string  `json:"message"`
}
type ImageResponse struct {
	Code    int    `json:"code"`
	Data    Image  `json:"data"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

//	type Images struct {
//		Imagelist []Image `json:"imagelist"`
//	}
type Image struct {
	Template_id     int    `json:"template_id"`
	Vm_info         []any  `json:"vm_info"`
	Image_type      string `json:"image_type"`
	Os_distribution string `json:"os_distribution"`
	Name            string `json:"name"`
	Image_id        string `json:"image_id"`
	Distro          string `json:"distro"`
	Sku_type        string `json:"sku_type"`
	Image_state     string `json:"image_state"`
	Running_vms     string `json:"running_vms"`
	Cloning_ops     string `json:"cloning_ops"`
	Image_size      string `json:"image_size"`
	Creation_time   string `json:"creation_time"`
}

type ImageDeleteBody struct {
	Action_type string `json:"action_type"`
}
