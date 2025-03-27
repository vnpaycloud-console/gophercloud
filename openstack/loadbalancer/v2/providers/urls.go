package providers

import "github.com/vnpaycloud-console/gophercloud/v2"

const (
	rootPath     = "lbaas"
	resourcePath = "providers"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}
