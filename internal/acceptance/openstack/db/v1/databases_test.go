//go:build acceptance || db || databases

package v1

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/db/v1/databases"
)

// Because it takes so long to create an instance,
// all tests will be housed in a single function.
func TestDatabases(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client, err := clients.NewDBV1Client()
	if err != nil {
		t.Fatalf("Unable to create a DB client: %v", err)
	}

	// Create and Get an instance.
	instance, err := CreateInstance(t, client)
	if err != nil {
		t.Fatalf("Unable to create instance: %v", err)
	}
	defer DeleteInstance(t, client, instance.ID)

	// Create a database.
	err = CreateDatabase(t, client, instance.ID)
	if err != nil {
		t.Fatalf("Unable to create database: %v", err)
	}

	// List all databases.
	allPages, err := databases.List(client, instance.ID).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list databases: %v", err)
	}

	allDatabases, err := databases.ExtractDBs(allPages)
	if err != nil {
		t.Fatalf("Unable to extract databases: %v", err)
	}

	for _, db := range allDatabases {
		tools.PrintResource(t, db)
	}

	defer DeleteDatabase(t, client, instance.ID, allDatabases[0].Name)

}
