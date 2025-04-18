package limits

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
)

// GetOptsBuilder allows extensions to add additional parameters to the
// Get request.
type GetOptsBuilder interface {
	ToLimitsQuery() (string, error)
}

// GetOpts enables retrieving limits by a specific tenant.
type GetOpts struct {
	// The tenant ID to retrieve limits for.
	TenantID string `q:"tenant_id"`
}

// ToLimitsQuery formats a GetOpts into a query string.
func (opts GetOpts) ToLimitsQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get returns the limits about the currently scoped tenant.
func Get(ctx context.Context, client *gophercloud.ServiceClient, opts GetOptsBuilder) (r GetResult) {
	url := getURL(client)
	if opts != nil {
		query, err := opts.ToLimitsQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	resp, err := client.Get(ctx, url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
