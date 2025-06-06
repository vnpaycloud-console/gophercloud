package testing

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/identity/v2/tokens"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
	"github.com/vnpaycloud-console/gophercloud/v2/testhelper/client"
)

func tokenPost(t *testing.T, options gophercloud.AuthOptions, requestJSON string) tokens.CreateResult {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleTokenPost(t, requestJSON)

	return tokens.Create(context.TODO(), client.ServiceClient(), options)
}

func tokenPostErr(t *testing.T, options gophercloud.AuthOptions, expectedErr error) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleTokenPost(t, "")

	actualErr := tokens.Create(context.TODO(), client.ServiceClient(), options).Err
	th.CheckDeepEquals(t, expectedErr, actualErr)
}

func TestCreateWithPassword(t *testing.T) {
	options := gophercloud.AuthOptions{
		Username: "me",
		Password: "swordfish",
	}

	IsSuccessful(t, tokenPost(t, options, `
    {
      "auth": {
        "passwordCredentials": {
          "username": "me",
          "password": "swordfish"
        }
      }
    }
  `))
}

func TestCreateTokenWithTenantID(t *testing.T) {
	options := gophercloud.AuthOptions{
		Username: "me",
		Password: "opensesame",
		TenantID: "fc394f2ab2df4114bde39905f800dc57",
	}

	IsSuccessful(t, tokenPost(t, options, `
    {
      "auth": {
        "tenantId": "fc394f2ab2df4114bde39905f800dc57",
        "passwordCredentials": {
          "username": "me",
          "password": "opensesame"
        }
      }
    }
  `))
}

func TestCreateTokenWithTenantName(t *testing.T) {
	options := gophercloud.AuthOptions{
		Username:   "me",
		Password:   "opensesame",
		TenantName: "demo",
	}

	IsSuccessful(t, tokenPost(t, options, `
    {
      "auth": {
        "tenantName": "demo",
        "passwordCredentials": {
          "username": "me",
          "password": "opensesame"
        }
      }
    }
  `))
}

func TestRequireUsername(t *testing.T) {
	options := gophercloud.AuthOptions{
		Password: "thing",
	}

	tokenPostErr(t, options, gophercloud.ErrMissingInput{Argument: "Username"})
}

func tokenGet(t *testing.T, tokenId string) tokens.GetResult {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleTokenGet(t, tokenId)
	return tokens.Get(context.TODO(), client.ServiceClient(), tokenId)
}

func TestGetWithToken(t *testing.T) {
	GetIsSuccessful(t, tokenGet(t, "db22caf43c934e6c829087c41ff8d8d6"))
}
