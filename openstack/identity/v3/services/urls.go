package services

import "github.com/vnpaycloud-console/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("services")
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("services")
}

func serviceURL(client *gophercloud.ServiceClient, serviceID string) string {
	return client.ServiceURL("services", serviceID)
}

func updateURL(client *gophercloud.ServiceClient, serviceID string) string {
	return client.ServiceURL("services", serviceID)
}
