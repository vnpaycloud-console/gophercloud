package tenants

import "github.com/vnpaycloud-console/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("tenants")
}

func getURL(client *gophercloud.ServiceClient, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("tenants")
}

func deleteURL(client *gophercloud.ServiceClient, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}

func updateURL(client *gophercloud.ServiceClient, tenantID string) string {
	return client.ServiceURL("tenants", tenantID)
}
