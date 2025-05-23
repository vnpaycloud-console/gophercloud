//go:build acceptance || db || users

package v1

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/db/v1/users"
)

// Because it takes so long to create an instance,
// all tests will be housed in a single function.
func TestUsers(t *testing.T) {
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

	// Create a user.
	err = CreateUser(t, client, instance.ID)
	if err != nil {
		t.Fatalf("Unable to create user: %v", err)
	}

	// List all users.
	allPages, err := users.List(client, instance.ID).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list users: %v", err)
	}

	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		t.Fatalf("Unable to extract users: %v", err)
	}

	for _, user := range allUsers {
		tools.PrintResource(t, user)
	}

	defer DeleteUser(t, client, instance.ID, allUsers[0].Name)
}
