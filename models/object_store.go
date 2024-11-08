package models

type ObjectStore struct {
	ID                           float64 `json:"id"`
	Name                         string  `json:"name"`
	Status                       string  `json:"status"`
	BucketSize                   string  `json:"bucket_size"`
	CreatedOn                    string  `json:"created_at"`
	VersioningStatus             string  `json:"versioning_status"`
	LifecycleConfigurationStatus string  `json:"lifecycle_configuration_status"`
}

type ObjectStorePayload struct {
	BucketName string `json:"bucket_name"`
	Region     string `json:"region"`
	ProjectID  int    `json:"project_id"`
}

type ResponseBuckets struct {
	Code    int           `json:"code"`
	Data    []ObjectStore `json:"data"`
	Error   string        `json:"error"`
	Message string        `json:"message"`
}
