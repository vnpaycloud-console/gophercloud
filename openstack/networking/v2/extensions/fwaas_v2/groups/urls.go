package groups

import "github.com/vnpaycloud-console/gophercloud/v2"

const (
	rootPath     = "fwaas"
	resourcePath = "firewall_groups"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}
