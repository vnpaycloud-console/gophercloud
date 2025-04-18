package swauth

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
)

// AuthOptsBuilder describes struct types that can be accepted by the Auth call.
type AuthOptsBuilder interface {
	ToAuthOptsMap() (map[string]string, error)
}

// AuthOpts specifies an authentication request.
type AuthOpts struct {
	// User is an Swauth-based username in username:tenant format.
	User string `h:"X-Auth-User" required:"true"`

	// Key is a secret/password to authenticate the User with.
	Key string `h:"X-Auth-Key" required:"true"`
}

// ToAuthOptsMap formats an AuthOpts structure into a request body.
func (opts AuthOpts) ToAuthOptsMap() (map[string]string, error) {
	return gophercloud.BuildHeaders(opts)
}

// Auth performs an authentication request for a Swauth-based user.
func Auth(ctx context.Context, c *gophercloud.ProviderClient, opts AuthOptsBuilder) (r GetAuthResult) {
	h := make(map[string]string)

	if opts != nil {
		headers, err := opts.ToAuthOptsMap()
		if err != nil {
			r.Err = err
			return
		}

		for k, v := range headers {
			h[k] = v
		}
	}

	resp, err := c.Request(ctx, "GET", getURL(c), &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return r
}

// NewObjectStorageV1 creates a Swauth-authenticated *gophercloud.ServiceClient
// client that can issue ObjectStorage-based API calls.
func NewObjectStorageV1(ctx context.Context, pc *gophercloud.ProviderClient, authOpts AuthOpts) (*gophercloud.ServiceClient, error) {
	auth, err := Auth(ctx, pc, authOpts).Extract()
	if err != nil {
		return nil, err
	}

	swiftClient := &gophercloud.ServiceClient{
		ProviderClient: pc,
		Endpoint:       gophercloud.NormalizeURL(auth.StorageURL),
	}

	swiftClient.TokenID = auth.Token

	return swiftClient, nil
}
