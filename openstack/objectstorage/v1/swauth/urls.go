package swauth

import "github.com/vnpaycloud-console/gophercloud/v2"

func getURL(c *gophercloud.ProviderClient) string {
	return c.IdentityBase + "auth/v1.0"
}
