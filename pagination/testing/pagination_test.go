package testing

import (
	"github.com/vnpaycloud-console/gophercloud/v2"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func createClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{TokenID: "abc123"},
		Endpoint:       th.Endpoint(),
	}
}
