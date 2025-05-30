//go:build acceptance || orchestration || stackresources

package v1

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/orchestration/v1/stackresources"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestStackResources(t *testing.T) {
	client, err := clients.NewOrchestrationV1Client()
	th.AssertNoErr(t, err)

	stack, err := CreateStack(t, client)
	th.AssertNoErr(t, err)
	defer DeleteStack(t, client, stack.Name, stack.ID)

	resource, err := stackresources.Get(context.TODO(), client, stack.Name, stack.ID, basicTemplateResourceName).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, resource)

	metadata, err := stackresources.Metadata(context.TODO(), client, stack.Name, stack.ID, basicTemplateResourceName).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, metadata)

	markUnhealthyOpts := &stackresources.MarkUnhealthyOpts{
		MarkUnhealthy:        true,
		ResourceStatusReason: "Wrong security policy is detected.",
	}

	err = stackresources.MarkUnhealthy(context.TODO(), client, stack.Name, stack.ID, basicTemplateResourceName, markUnhealthyOpts).ExtractErr()
	th.AssertNoErr(t, err)

	unhealthyResource, err := stackresources.Get(context.TODO(), client, stack.Name, stack.ID, basicTemplateResourceName).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "CHECK_FAILED", unhealthyResource.Status)
	tools.PrintResource(t, unhealthyResource)

	allPages, err := stackresources.List(client, stack.Name, stack.ID, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allResources, err := stackresources.ExtractResources(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, v := range allResources {
		if v.Name == basicTemplateResourceName {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}
