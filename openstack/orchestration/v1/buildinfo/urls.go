package buildinfo

import "github.com/vnpaycloud-console/gophercloud/v2"

func getURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("build_info")
}
