package clusters

import (
	"github.com/vnpaycloud-console/gophercloud/v2"
)

var apiName = "clusters"

func commonURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiName)
}

func idURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL(apiName, id)
}

func createURL(client *gophercloud.ServiceClient) string {
	return commonURL(client)
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return idURL(client, id)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("clusters", id)
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("clusters")
}

func listDetailURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("clusters", "detail")
}

func updateURL(client *gophercloud.ServiceClient, id string) string {
	return idURL(client, id)
}

func upgradeURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("clusters", id, "actions/upgrade")
}

func resizeURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("clusters", id, "actions/resize")
}
