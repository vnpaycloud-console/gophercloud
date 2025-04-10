//go:build acceptance || identity || services

package v3

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/identity/v3/services"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestServicesList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	listOpts := services.ListOpts{
		ServiceType: "identity",
	}

	allPages, err := services.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServices, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, service := range allServices {
		tools.PrintResource(t, service)

		if service.Type == "identity" {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestServicesCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	createOpts := services.CreateOpts{
		Type: "testing",
		Extra: map[string]any{
			"email": "testservice@example.com",
		},
	}

	// Create service in the default domain
	service, err := CreateService(t, client, &createOpts)
	th.AssertNoErr(t, err)
	defer DeleteService(t, client, service.ID)

	tools.PrintResource(t, service)
	tools.PrintResource(t, service.Extra)

	updateOpts := services.UpdateOpts{
		Type: "testing2",
		Extra: map[string]any{
			"description": "Test Users",
			"email":       "thetestservice@example.com",
		},
	}

	newService, err := services.Update(context.TODO(), client, service.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newService)
	tools.PrintResource(t, newService.Extra)

	th.AssertEquals(t, newService.Extra["description"], "Test Users")
}
