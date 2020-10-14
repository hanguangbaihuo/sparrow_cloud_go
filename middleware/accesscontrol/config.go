package accesscontrol

// Config is a struct for specifying configuration options for the accesscontrol middleware.
type Config struct {
	// AccessControlService is access control service address:port
	AccessControlService string
	// APIPath
	APIPath string
	// ServiceName is the service name, used for app_name
	ServiceName string
	// SkipAccessContorl, true: skip the access control; false: check user's resource
	SkipAccessContorl bool
}

// ACRequestData is request access control service data
type ACRequestData struct {
	UserID       string `json:"user_id"`
	AppName      string `json:"app_name"`
	ResourceCode string `json:"resource_code"`
}

// ACResponse is access control response when reqeust ac service
type ACResponse struct {
	HasPerm bool `json:"has_perm"`
}
