package apiversions

import (
	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// ListVersions lists all the Neutron API versions available to end-users.
func ListVersions(c *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, listURL(c), func(r pagination.PageResult) pagination.Page {
		return APIVersionPage{pagination.SinglePageBase(r)}
	})
}

// ListVersionResources lists all of the different API resources for a
// particular API versions. Typical resources for Neutron might be: networks,
// subnets, etc.
func ListVersionResources(c *gophercloud.ServiceClient, v string) pagination.Pager {
	return pagination.NewPager(c, getURL(c, v), func(r pagination.PageResult) pagination.Page {
		return APIVersionResourcePage{pagination.SinglePageBase(r)}
	})
}
