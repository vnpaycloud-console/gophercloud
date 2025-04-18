package trunks

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/networking/v2/extensions/trunks"
)

func CreateTrunk(t *testing.T, client *gophercloud.ServiceClient, parentPortID string, subportIDs ...string) (trunk *trunks.Trunk, err error) {
	trunkName := tools.RandomString("TESTACC-", 8)
	iTrue := true
	opts := trunks.CreateOpts{
		Name:         trunkName,
		Description:  "Trunk created by gophercloud",
		AdminStateUp: &iTrue,
		PortID:       parentPortID,
	}

	opts.Subports = make([]trunks.Subport, len(subportIDs))
	for id, subportID := range subportIDs {
		opts.Subports[id] = trunks.Subport{
			SegmentationID:   id + 1,
			SegmentationType: "vlan",
			PortID:           subportID,
		}
	}

	t.Logf("Attempting to create trunk: %s", opts.Name)
	trunk, err = trunks.Create(context.TODO(), client, opts).Extract()
	if err == nil {
		t.Logf("Successfully created trunk")
	}
	return
}

func DeleteTrunk(t *testing.T, client *gophercloud.ServiceClient, trunkID string) {
	t.Logf("Attempting to delete trunk: %s", trunkID)
	err := trunks.Delete(context.TODO(), client, trunkID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete trunk %s: %v", trunkID, err)
	}

	t.Logf("Deleted trunk: %s", trunkID)
}
