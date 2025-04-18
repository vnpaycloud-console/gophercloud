package testing

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/openstack/compute/v2/diagnostics"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
	"github.com/vnpaycloud-console/gophercloud/v2/testhelper/client"
)

func TestGetDiagnostics(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleDiagnosticGetSuccessfully(t)

	expected := map[string]any{"cpu0_time": float64(173), "memory": float64(524288)}

	res, err := diagnostics.Get(context.TODO(), client.ServiceClient(), "1234asdf").Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, expected, res)
}
