package drivers

import "github.com/vnpaycloud-console/gophercloud/v2"

func driversURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("drivers")
}

func driverDetailsURL(client *gophercloud.ServiceClient, driverName string) string {
	return client.ServiceURL("drivers", driverName)
}

func driverPropertiesURL(client *gophercloud.ServiceClient, driverName string) string {
	return client.ServiceURL("drivers", driverName, "properties")
}

func driverDiskPropertiesURL(client *gophercloud.ServiceClient, driverName string) string {
	return client.ServiceURL("drivers", driverName, "raid", "logical_disk_properties")
}
