//go:build acceptance || identity || groups

package v3

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/identity/v3/groups"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestGroupCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	description := "Test Groups"
	domainID := "default"
	createOpts := groups.CreateOpts{
		Description: description,
		DomainID:    domainID,
		Extra: map[string]any{
			"email": "testgroup@example.com",
		},
	}

	// Create Group in the default domain
	group, err := CreateGroup(t, client, &createOpts)
	th.AssertNoErr(t, err)
	defer DeleteGroup(t, client, group.ID)

	tools.PrintResource(t, group)
	tools.PrintResource(t, group.Extra)

	th.AssertEquals(t, group.Description, description)
	th.AssertEquals(t, group.DomainID, domainID)
	th.AssertDeepEquals(t, group.Extra, createOpts.Extra)

	description = ""
	updateOpts := groups.UpdateOpts{
		Description: &description,
		Extra: map[string]any{
			"email": "thetestgroup@example.com",
		},
	}

	newGroup, err := groups.Update(context.TODO(), client, group.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newGroup)
	tools.PrintResource(t, newGroup.Extra)

	th.AssertEquals(t, newGroup.Description, description)
	th.AssertDeepEquals(t, newGroup.Extra, updateOpts.Extra)

	listOpts := groups.ListOpts{
		DomainID: "default",
	}

	// List all Groups in default domain
	allPages, err := groups.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allGroups, err := groups.ExtractGroups(allPages)
	th.AssertNoErr(t, err)

	for _, g := range allGroups {
		tools.PrintResource(t, g)
		tools.PrintResource(t, g.Extra)
	}

	var found bool
	for _, group := range allGroups {
		tools.PrintResource(t, group)
		tools.PrintResource(t, group.Extra)

		if group.Name == newGroup.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	listOpts.Filters = map[string]string{
		"name__contains": "TEST",
	}

	allPages, err = groups.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allGroups, err = groups.ExtractGroups(allPages)
	th.AssertNoErr(t, err)

	found = false
	for _, group := range allGroups {
		tools.PrintResource(t, group)
		tools.PrintResource(t, group.Extra)

		if group.Name == newGroup.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	listOpts.Filters = map[string]string{
		"name__contains": "foo",
	}

	allPages, err = groups.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allGroups, err = groups.ExtractGroups(allPages)
	th.AssertNoErr(t, err)

	found = false
	for _, group := range allGroups {
		tools.PrintResource(t, group)
		tools.PrintResource(t, group.Extra)

		if group.Name == newGroup.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, false)

	// Get the recently created group by ID
	p, err := groups.Get(context.TODO(), client, group.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, p)
}
