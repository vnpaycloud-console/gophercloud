package extensions

import "github.com/vnpaycloud-console/gophercloud/v2"

func ActionURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "action")
}
