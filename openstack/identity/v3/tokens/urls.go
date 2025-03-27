package tokens

import "github.com/vnpaycloud-console/gophercloud/v2"

func tokenURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
