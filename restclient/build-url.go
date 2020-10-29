package restclient

import (
	"fmt"
	"log"
	"net/url"
)

func buildURL(protocol string, address string, apiPath string) string {
	tempURL := fmt.Sprintf("%s://%s%s", protocol, address, apiPath)
	destURL, err := url.Parse(tempURL)
	if err != nil {
		log.Printf("parse url: %s occur error: %s\n", tempURL, err)
		return tempURL
	}
	destURL.RawQuery = destURL.Query().Encode()
	return destURL.String()
}
