package testing

import (
	"context"
	"testing"
	"time"

	"github.com/vnpaycloud-console/gophercloud/v2/openstack/sharedfilesystems/v2/availabilityzones"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
	"github.com/vnpaycloud-console/gophercloud/v2/testhelper/client"
)

// Verifies that availability zones can be listed correctly
func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	allPages, err := availabilityzones.List(client.ServiceClient()).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := availabilityzones.ExtractAvailabilityZones(allPages)
	th.AssertNoErr(t, err)
	var nilTime time.Time
	expected := []availabilityzones.AvailabilityZone{
		{
			Name:      "nova",
			CreatedAt: time.Date(2015, 9, 18, 9, 50, 55, 0, time.UTC),
			UpdatedAt: nilTime,
			ID:        "388c983d-258e-4a0e-b1ba-10da37d766db",
		},
	}

	th.CheckDeepEquals(t, expected, actual)
}
