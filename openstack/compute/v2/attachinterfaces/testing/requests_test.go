package testing

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/openstack/compute/v2/attachinterfaces"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
	"github.com/vnpaycloud-console/gophercloud/v2/testhelper/client"
)

func TestListInterface(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleInterfaceListSuccessfully(t)

	expected := ListInterfacesExpected
	pages := 0
	err := attachinterfaces.List(client.ServiceClient(), "b07e7a3b-d951-4efc-a4f9-ac9f001afb7f").EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := attachinterfaces.ExtractInterfaces(page)
		th.AssertNoErr(t, err)

		if len(actual) != 1 {
			t.Fatalf("Expected 1 interface, got %d", len(actual))
		}
		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, pages)
}

func TestListInterfacesAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleInterfaceListSuccessfully(t)

	allPages, err := attachinterfaces.List(client.ServiceClient(), "b07e7a3b-d951-4efc-a4f9-ac9f001afb7f").AllPages(context.TODO())
	th.AssertNoErr(t, err)
	_, err = attachinterfaces.ExtractInterfaces(allPages)
	th.AssertNoErr(t, err)
}

func TestGetInterface(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleInterfaceGetSuccessfully(t)

	expected := GetInterfaceExpected

	serverID := "b07e7a3b-d951-4efc-a4f9-ac9f001afb7f"
	interfaceID := "0dde1598-b374-474e-986f-5b8dd1df1d4e"

	actual, err := attachinterfaces.Get(context.TODO(), client.ServiceClient(), serverID, interfaceID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expected, actual)
}

func TestCreateInterface(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleInterfaceCreateSuccessfully(t)

	expected := CreateInterfacesExpected

	serverID := "b07e7a3b-d951-4efc-a4f9-ac9f001afb7f"
	networkID := "8a5fe506-7e9f-4091-899b-96336909d93c"

	actual, err := attachinterfaces.Create(context.TODO(), client.ServiceClient(), serverID, attachinterfaces.CreateOpts{
		NetworkID: networkID,
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expected, actual)
}

func TestDeleteInterface(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleInterfaceDeleteSuccessfully(t)

	serverID := "b07e7a3b-d951-4efc-a4f9-ac9f001afb7f"
	portID := "0dde1598-b374-474e-986f-5b8dd1df1d4e"

	err := attachinterfaces.Delete(context.TODO(), client.ServiceClient(), serverID, portID).ExtractErr()
	th.AssertNoErr(t, err)
}
