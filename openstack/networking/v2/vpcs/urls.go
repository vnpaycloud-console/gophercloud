package vpcs

import "github.com/vnpaycloud-console/gophercloud/v2"

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("vpcs", id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("vpcs")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}
