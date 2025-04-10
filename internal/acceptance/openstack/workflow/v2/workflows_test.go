//go:build acceptance || workflow || workflows

package v2

import (
	"testing"
	"time"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/workflow/v2/workflows"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestWorkflowsCreateGetDelete(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)

	workflow, err := CreateWorkflow(t, client)
	th.AssertNoErr(t, err)
	defer DeleteWorkflow(t, client, workflow)

	workflowget, err := GetWorkflow(t, client, workflow.ID)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, workflowget)
}

func TestWorkflowsList(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)
	workflow, err := CreateWorkflow(t, client)
	th.AssertNoErr(t, err)
	defer DeleteWorkflow(t, client, workflow)
	list, err := ListWorkflows(t, client, &workflows.ListOpts{
		Name: &workflows.ListFilter{
			Value: workflow.Name,
		},
		Tags: []string{"tag1"},
		CreatedAt: &workflows.ListDateFilter{
			Filter: workflows.FilterGT,
			Value:  time.Now().AddDate(-1, 0, 0),
		},
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(list))
	tools.PrintResource(t, list)
}
