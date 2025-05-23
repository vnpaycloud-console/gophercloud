package attachments

import "github.com/vnpaycloud-console/gophercloud/v2"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("attachments")
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("attachments", "detail")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("attachments", id)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("attachments", id)
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("attachments", id)
}

func completeURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("attachments", id, "action")
}
