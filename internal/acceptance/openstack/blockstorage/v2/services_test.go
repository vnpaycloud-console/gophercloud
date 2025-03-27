//go:build acceptance || blockstorage || services

package v2

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/blockstorage/v2/services"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestServicesList(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")
	clients.RequireAdmin(t)

	blockClient, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	allPages, err := services.List(blockClient, services.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServices, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	for _, service := range allServices {
		tools.PrintResource(t, service)
	}
}
