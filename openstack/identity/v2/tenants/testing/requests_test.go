package testing

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/identity/v2/tenants"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
	"github.com/vnpaycloud-console/gophercloud/v2/testhelper/client"
)

func TestListTenants(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTenantsSuccessfully(t)

	count := 0
	err := tenants.List(client.ServiceClient(), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := tenants.ExtractTenants(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedTenantSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestCreateTenant(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockCreateTenantResponse(t)

	opts := tenants.CreateOpts{
		Name:        "new_tenant",
		Description: "This is new tenant",
		Enabled:     gophercloud.Enabled,
	}

	tenant, err := tenants.Create(context.TODO(), client.ServiceClient(), opts).Extract()

	th.AssertNoErr(t, err)

	expected := &tenants.Tenant{
		Name:        "new_tenant",
		Description: "This is new tenant",
		Enabled:     true,
		ID:          "5c62ef576dc7444cbb73b1fe84b97648",
	}

	th.AssertDeepEquals(t, expected, tenant)
}

func TestDeleteTenant(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockDeleteTenantResponse(t)

	err := tenants.Delete(context.TODO(), client.ServiceClient(), "2466f69cd4714d89a548a68ed97ffcd4").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUpdateTenant(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockUpdateTenantResponse(t)

	id := "5c62ef576dc7444cbb73b1fe84b97648"
	description := "This is new name"
	opts := tenants.UpdateOpts{
		Name:        "new_name",
		Description: &description,
		Enabled:     gophercloud.Enabled,
	}

	tenant, err := tenants.Update(context.TODO(), client.ServiceClient(), id, opts).Extract()

	th.AssertNoErr(t, err)

	expected := &tenants.Tenant{
		Name:        "new_name",
		ID:          id,
		Description: "This is new name",
		Enabled:     true,
	}

	th.AssertDeepEquals(t, expected, tenant)
}

func TestGetTenant(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockGetTenantResponse(t)

	tenant, err := tenants.Get(context.TODO(), client.ServiceClient(), "5c62ef576dc7444cbb73b1fe84b97648").Extract()
	th.AssertNoErr(t, err)

	expected := &tenants.Tenant{
		Name:        "new_tenant",
		ID:          "5c62ef576dc7444cbb73b1fe84b97648",
		Description: "This is new tenant",
		Enabled:     true,
	}

	th.AssertDeepEquals(t, expected, tenant)
}
