//go:build acceptance || orchestration || stacks

package v1

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/orchestration/v1/stacks"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestStacksCRUD(t *testing.T) {
	client, err := clients.NewOrchestrationV1Client()
	th.AssertNoErr(t, err)

	createdStack, err := CreateStack(t, client)
	th.AssertNoErr(t, err)
	defer DeleteStack(t, client, createdStack.Name, createdStack.ID)

	tools.PrintResource(t, createdStack)
	tools.PrintResource(t, createdStack.CreationTime)

	template := new(stacks.Template)
	template.Bin = []byte(basicTemplate)
	updateOpts := stacks.UpdateOpts{
		TemplateOpts: template,
		Timeout:      20,
	}

	err = stacks.Update(context.TODO(), client, createdStack.Name, createdStack.ID, updateOpts).ExtractErr()
	th.AssertNoErr(t, err)

	err = WaitForStackStatus(client, createdStack.Name, createdStack.ID, "UPDATE_COMPLETE")
	th.AssertNoErr(t, err)

	var found bool
	allPages, err := stacks.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allStacks, err := stacks.ExtractStacks(allPages)
	th.AssertNoErr(t, err)

	for _, v := range allStacks {
		if v.ID == createdStack.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}
