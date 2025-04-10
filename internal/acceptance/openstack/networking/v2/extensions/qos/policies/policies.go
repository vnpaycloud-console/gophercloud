package policies

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/networking/v2/extensions/qos/policies"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

// CreateQoSPolicy will create a QoS policy. An error will be returned if the
// QoS policy could not be created.
func CreateQoSPolicy(t *testing.T, client *gophercloud.ServiceClient) (*policies.Policy, error) {
	policyName := tools.RandomString("TESTACC-", 8)
	policyDescription := tools.RandomString("TESTACC-DESC-", 8)

	createOpts := policies.CreateOpts{
		Name:        policyName,
		Description: policyDescription,
	}

	t.Logf("Attempting to create a QoS policy: %s", policyName)

	policy, err := policies.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Succesfully created a QoS policy")

	th.AssertEquals(t, policyName, policy.Name)
	th.AssertEquals(t, policyDescription, policy.Description)

	return policy, nil
}

// DeleteQoSPolicy will delete a QoS policy with a specified ID.
// A fatal error will occur if the delete was not successful.
func DeleteQoSPolicy(t *testing.T, client *gophercloud.ServiceClient, policyID string) {
	t.Logf("Attempting to delete the QoS policy: %s", policyID)

	err := policies.Delete(context.TODO(), client, policyID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete QoS policy %s: %v", policyID, err)
	}

	t.Logf("Deleted QoS policy: %s", policyID)
}
