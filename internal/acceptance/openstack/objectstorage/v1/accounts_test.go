//go:build acceptance || objectstorage || accounts

package v1

import (
	"context"
	"strings"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/objectstorage/v1/accounts"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestAccounts(t *testing.T) {
	client, err := clients.NewObjectStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	// Update an account's metadata.
	metadata := map[string]string{
		"Gophercloud-Test": "accounts",
	}
	updateres := accounts.Update(context.TODO(), client, accounts.UpdateOpts{Metadata: metadata})
	t.Logf("Update Account Response: %+v\n", updateres)
	updateHeaders, err := updateres.Extract()
	th.AssertNoErr(t, err)
	t.Logf("Update Account Response Headers: %+v\n", updateHeaders)

	// Defer the deletion of the metadata set above.
	defer func() {
		tempMap := make(map[string]string)
		for k := range metadata {
			tempMap[k] = ""
		}
		updateres = accounts.Update(context.TODO(), client, accounts.UpdateOpts{Metadata: tempMap})
		th.AssertNoErr(t, updateres.Err)
	}()

	// Extract the custom metadata from the 'Get' response.
	res := accounts.Get(context.TODO(), client, nil)

	h, err := res.Extract()
	th.AssertNoErr(t, err)
	t.Logf("Get Account Response Headers: %+v\n", h)

	am, err := res.ExtractMetadata()
	th.AssertNoErr(t, err)
	for k := range metadata {
		if am[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
			return
		}
	}
}
