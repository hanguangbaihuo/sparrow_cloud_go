package restclient

import (
	"fmt"
)

func buildURL(protocol string, address string, apiPath string) string {
	return fmt.Sprintf("%s://%s%s", protocol, address, apiPath)
}
