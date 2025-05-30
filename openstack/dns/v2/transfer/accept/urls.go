package accept

import "github.com/vnpaycloud-console/gophercloud/v2"

const (
	rootPath     = "zones"
	tasksPath    = "tasks"
	resourcePath = "transfer_accepts"
)

func baseURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, tasksPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, transferAcceptID string) string {
	return c.ServiceURL(rootPath, tasksPath, resourcePath, transferAcceptID)
}
