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

var apiParam = "?user_id=%s&app_name=%s&resource_code=%s"

// ACResponse is access control response when reqeust ac service
type ACResponse struct {
	HasPerm bool `json:"has_perm"`
}
