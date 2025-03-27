package portgroups

import "github.com/vnpaycloud-console/gophercloud/v2"

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("portgroups")
}

func listURL(client *gophercloud.ServiceClient) string {
	return createURL(client)
}

func resourceURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("portgroups", id)
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return resourceURL(client, id)
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return resourceURL(client, id)
}
