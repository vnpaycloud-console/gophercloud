package datastores

import "github.com/vnpaycloud-console/gophercloud/v2"

func baseURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("datastores")
}

func resourceURL(c *gophercloud.ServiceClient, dsID string) string {
	return c.ServiceURL("datastores", dsID)
}

func versionsURL(c *gophercloud.ServiceClient, dsID string) string {
	return c.ServiceURL("datastores", dsID, "versions")
}

func versionURL(c *gophercloud.ServiceClient, dsID, versionID string) string {
	return c.ServiceURL("datastores", dsID, "versions", versionID)
}
