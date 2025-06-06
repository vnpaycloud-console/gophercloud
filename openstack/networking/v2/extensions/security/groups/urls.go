package groups

import "github.com/vnpaycloud-console/gophercloud/v2"

const rootPath = "security-groups"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id)
}
