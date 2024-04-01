package models

type SshKeyResponse struct {
	Code    int           `json:"code"`
	Data    []SshKey      `json:"data"`
	Error   []interface{} `json:"error"`
	Message string        `json:"message"`
}
type SshKey struct {
	Label     string `json:"label"`
	Ssh_key   string `json:"ssh_key"`
	Pk        int    `json:"pk"`
	Timestamp string `json:"timestamp"`
}

type AddSshKey struct {
	Label  string `json:"label"`
	SshKey string `json:"ssh_key"`
}
