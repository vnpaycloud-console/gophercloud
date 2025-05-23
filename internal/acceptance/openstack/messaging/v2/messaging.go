package v2

import (
	"context"
	"strings"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/messaging/v2/claims"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/messaging/v2/messages"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/messaging/v2/queues"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func CreateQueue(t *testing.T, client *gophercloud.ServiceClient) (string, error) {
	queueName := tools.RandomString("ACPTTEST", 5)

	t.Logf("Attempting to create Queue: %s", queueName)

	createOpts := queues.CreateOpts{
		QueueName:                  queueName,
		MaxMessagesPostSize:        262143,
		DefaultMessageTTL:          3700,
		DeadLetterQueueMessagesTTL: 3500,
		MaxClaimCount:              10,
		Extra:                      map[string]any{"description": "Test Queue for Gophercloud acceptance tests."},
	}

	createErr := queues.Create(context.TODO(), client, createOpts).ExtractErr()
	if createErr != nil {
		t.Fatalf("Unable to create Queue: %v", createErr)
	}

	_, err := GetQueue(t, client, queueName)
	th.AssertNoErr(t, err)

	t.Logf("Created Queue: %s", queueName)
	return queueName, nil
}

func DeleteQueue(t *testing.T, client *gophercloud.ServiceClient, queueName string) {
	t.Logf("Attempting to delete Queue: %s", queueName)
	err := queues.Delete(context.TODO(), client, queueName).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete Queue %s: %v", queueName, err)
	}

	t.Logf("Deleted Queue: %s", queueName)
}

func GetQueue(t *testing.T, client *gophercloud.ServiceClient, queueName string) (queues.QueueDetails, error) {
	t.Logf("Attempting to get Queue: %s", queueName)
	queue, err := queues.Get(context.TODO(), client, queueName).Extract()
	if err != nil {
		t.Fatalf("Unable to get Queue %s: %v", queueName, err)
	}
	return queue, nil
}

func CreateShare(t *testing.T, client *gophercloud.ServiceClient, queueName string) (queues.QueueShare, error) {
	t.Logf("Attempting to create share for queue: %s", queueName)

	shareOpts := queues.ShareOpts{
		Paths:   []queues.SharePath{queues.PathMessages},
		Methods: []queues.ShareMethod{queues.MethodPost},
	}

	share, err := queues.Share(context.TODO(), client, queueName, shareOpts).Extract()

	return share, err
}

func CreateMessage(t *testing.T, client *gophercloud.ServiceClient, queueName string) (messages.ResourceList, error) {
	t.Logf("Attempting to add message to Queue: %s", queueName)
	createOpts := messages.BatchCreateOpts{
		messages.CreateOpts{
			TTL:  300,
			Body: map[string]any{"Key": tools.RandomString("ACPTTEST", 8)},
		},
	}

	resource, err := messages.Create(context.TODO(), client, queueName, createOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to add message to queue %s: %v", queueName, err)
	} else {
		t.Logf("Successfully added message to queue: %s", queueName)
	}

	return resource, err
}

func ListMessages(t *testing.T, client *gophercloud.ServiceClient, queueName string) ([]messages.Message, error) {
	listOpts := messages.ListOpts{}
	var allMessages []messages.Message
	var listErr error

	t.Logf("Attempting to list messages on queue: %s", queueName)
	pager := messages.List(client, queueName, listOpts)
	err := pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		allMessages, listErr = messages.ExtractMessages(page)
		if listErr != nil {
			t.Fatalf("Unable to extract messages: %v", listErr)
		}

		for _, message := range allMessages {
			tools.PrintResource(t, message)
		}

		return true, nil
	})
	return allMessages, err
}

func CreateClaim(t *testing.T, client *gophercloud.ServiceClient, queueName string) ([]claims.Messages, error) {
	createOpts := claims.CreateOpts{}

	t.Logf("Attempting to create claim on queue: %s", queueName)
	claimedMessages, err := claims.Create(context.TODO(), client, queueName, createOpts).Extract()
	tools.PrintResource(t, claimedMessages)
	if err != nil {
		t.Fatalf("Unable to create claim: %v", err)
	}

	return claimedMessages, err
}

func GetClaim(t *testing.T, client *gophercloud.ServiceClient, queueName string, claimID string) (*claims.Claim, error) {
	t.Logf("Attempting to get claim: %s", claimID)
	claim, err := claims.Get(context.TODO(), client, queueName, claimID).Extract()
	if err != nil {
		t.Fatalf("Unable to get claim: %s", claimID)
	}

	return claim, err
}

func DeleteClaim(t *testing.T, client *gophercloud.ServiceClient, queueName string, claimID string) error {
	t.Logf("Attempting to delete claim: %s", claimID)
	err := claims.Delete(context.TODO(), client, queueName, claimID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete claim: %s", claimID)
	}
	t.Logf("Sucessfully deleted claim: %s", claimID)

	return err
}

func ExtractIDs(claim []claims.Messages) ([]string, []string) {
	var claimIDs []string
	var messageID []string

	for _, msg := range claim {
		parts := strings.Split(msg.Href, "?claim_id=")
		if len(parts) == 2 {
			pieces := strings.Split(parts[0], "/")
			if len(pieces) > 0 {
				messageID = append(messageID, pieces[len(pieces)-1])
			}
			claimIDs = append(claimIDs, parts[1])
		}
	}
	encountered := map[string]bool{}
	for v := range claimIDs {
		encountered[claimIDs[v]] = true
	}

	var uniqueClaimIDs []string

	for key := range encountered {
		uniqueClaimIDs = append(uniqueClaimIDs, key)
	}
	return uniqueClaimIDs, messageID
}
