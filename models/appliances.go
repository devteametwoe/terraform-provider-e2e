package models

type LoadBalancerCreate struct {
	PlanName         string                 `json:"plan_name"`
	LbName           string                 `json:"lb_name"`
	LbType           string                 `json:"lb_type,omitempty"`
	LbMode           string                 `json:"lb_mode"`
	LbPort           string                 `json:"lb_port"`
	NodeListType     string                 `json:"node_list_type"`
	CheckBoxEnable   string                 `json:"checkbox_enable"`
	LbReserveIp      string                 `json:"lb_reserve_ip"`
	SslCertificateId string                 `json:"ssl_certificate_id"`
	SslContext       map[string]interface{} `json:"ssl_context"`
	EnableBitninja   bool                   `json:"enable_bitninja"`
	Backends         []Backend              `json:"backends"`
	AclList          []AclListInfo          `json:"acl_list"`
	AclMap           []AclMapInfo           `json:"acl_map"`
	VpcList          []VpcDetail            `json:"vpc_list"`
	EnableEosLogger  EosDetail              `json:"enable_eos_logger,omitempty"`
	TcpBackend       []TcpBackendDetail     `json:"tcp_backend"`
	IsIpv6Attached   bool                   `json:"is_ipv6_attached"`
	DefaultBackend   string                 `json:"default_backend"`
}

type Backend struct {
	Balance        string   `json:"balance"`
	CheckboxEnable bool     `json:"checkbox_enable"`
	DomainName     string   `json:"domain_name"`
	CheckUrl       string   `json:"check_url"`
	Servers        []Server `json:"servers"`
	HttpCheck      bool     `json:"http_check"`
	Name           string   `json:"name"`
	ScalerId       string   `json:"scaler_id"`
	ScalerPort     string   `json:"scaler_port"`
}

type Server struct {
	BackendName string `json:"backend_name"`
	BackendIp   string `json:"backend_ip"`
	BackendPort string `json:"backend_port"`
}

type EosDetail struct {
	ApplianceId int    `json:"appliance_id"`
	AccessKey   string `json:"access_key"`
	Secretkey   string `json:"secret_key"`
	Bucket      string `json:"bucket"`
}

type TcpBackendDetail struct {
	BackendName string   `json:"backend_name"`
	Port        string   `json:"port"`
	Balance     string   `json:"balance"`
	Servers     []Server `json:"servers"`
}

type AclListInfo struct {
	AclName         string `json:"acl_name,omitempty"`
	AclCondition    string `json:"acl_condition,omitempty"`
	AclMatchingPath string `json:"acl_matching_path,omitempty"`
}

type AclMapInfo struct {
	AclBackend        string `json:"acl_backend,omitempty"`
	AclConditionState bool   `json:"acl_condition_state,omitempty"`
	AclName           string `json:"acl_name,omitempty"`
}

type VpcDetail struct {
	VpcName    string  `json:"vpc_name,omitempty"`
	Ipv4_cidr  string  `json:"ipv4_cidr,omitempty"`
	Network_id float64 `json:"network_id,omitempty"`
}
