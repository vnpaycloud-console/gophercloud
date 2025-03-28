package testing

import (
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/openstack/baremetalintrospection/noauth"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestNoAuth(t *testing.T) {
	noauthClient, err := noauth.NewBareMetalIntrospectionNoAuth(noauth.EndpointOpts{
		IronicInspectorEndpoint: "http://ironic:5050/v1",
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "", noauthClient.TokenID)
}
