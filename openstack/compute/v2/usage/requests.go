package usage

import (
	"net/url"
	"time"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// SingleTenantOpts are options for fetching usage of a single tenant.
type SingleTenantOpts struct {
	// The ending time to calculate usage statistics on compute and storage resources.
	End *time.Time `q:"end"`

	// The beginning time to calculate usage statistics on compute and storage resources.
	Start *time.Time `q:"start"`

	// Limit limits the amount of results returned by the API.
	// This requires the client to be set to microversion 2.40 or later.
	Limit int `q:"limit"`

	// Marker instructs the API call where to start listing from.
	// This requires the client to be set to microversion 2.40 or later.
	Marker string `q:"marker"`
}

// SingleTenantOptsBuilder allows extensions to add additional parameters to the
// SingleTenant request.
type SingleTenantOptsBuilder interface {
	ToUsageSingleTenantQuery() (string, error)
}

// ToUsageSingleTenantQuery formats a SingleTenantOpts into a query string.
func (opts SingleTenantOpts) ToUsageSingleTenantQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}

	params := q.Query()

	if opts.Start != nil {
		params.Add("start", opts.Start.Format(gophercloud.RFC3339MilliNoZ))
	}

	if opts.End != nil {
		params.Add("end", opts.End.Format(gophercloud.RFC3339MilliNoZ))
	}

	q = &url.URL{RawQuery: params.Encode()}
	return q.String(), nil
}

// SingleTenant returns usage data about a single tenant.
func SingleTenant(client *gophercloud.ServiceClient, tenantID string, opts SingleTenantOptsBuilder) pagination.Pager {
	url := getTenantURL(client, tenantID)
	if opts != nil {
		query, err := opts.ToUsageSingleTenantQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SingleTenantPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// AllTenantsOpts are options for fetching usage of all tenants.
type AllTenantsOpts struct {
	// Detailed will return detailed results.
	Detailed bool

	// The ending time to calculate usage statistics on compute and storage resources.
	End *time.Time `q:"end"`

	// The beginning time to calculate usage statistics on compute and storage resources.
	Start *time.Time `q:"start"`

	// Limit limits the amount of results returned by the API.
	// This requires the client to be set to microversion 2.40 or later.
	Limit int `q:"limit"`

	// Marker instructs the API call where to start listing from.
	// This requires the client to be set to microversion 2.40 or later.
	Marker string `q:"marker"`
}

// AllTenantsOptsBuilder allows extensions to add additional parameters to the
// AllTenants request.
type AllTenantsOptsBuilder interface {
	ToUsageAllTenantsQuery() (string, error)
}

// ToUsageAllTenantsQuery formats a AllTenantsOpts into a query string.
func (opts AllTenantsOpts) ToUsageAllTenantsQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}

	params := q.Query()

	if opts.Start != nil {
		params.Add("start", opts.Start.Format(gophercloud.RFC3339MilliNoZ))
	}

	if opts.End != nil {
		params.Add("end", opts.End.Format(gophercloud.RFC3339MilliNoZ))
	}

	if opts.Detailed {
		params.Add("detailed", "1")
	}

	q = &url.URL{RawQuery: params.Encode()}
	return q.String(), nil
}

// AllTenants returns usage data about all tenants.
func AllTenants(client *gophercloud.ServiceClient, opts AllTenantsOptsBuilder) pagination.Pager {
	url := allTenantsURL(client)
	if opts != nil {
		query, err := opts.ToUsageAllTenantsQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AllTenantsPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
