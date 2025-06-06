package capsules

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToCapsuleCreateMap() (map[string]any, error)
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToCapsuleListQuery() (string, error)
}

// Get requests details on a single capsule, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 203},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// A structure that contains either the template file or url. Call the
	// associated methods to extract the information relevant to send in a create request.
	TemplateOpts *Template `json:"-" required:"true"`
}

// ToCapsuleCreateMap assembles a request body based on the contents of
// a CreateOpts.
func (opts CreateOpts) ToCapsuleCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if err := opts.TemplateOpts.Parse(); err != nil {
		return nil, err
	}
	b["template"] = string(opts.TemplateOpts.Bin)

	return b, nil
}

// Create implements create capsule request.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCapsuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the capsule attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListOpts struct {
	Marker      string `q:"marker"`
	Limit       int    `q:"limit"`
	SortKey     string `q:"sort_key"`
	SortDir     string `q:"sort_dir"`
	AllProjects bool   `q:"all_projects"`
}

// ToCapsuleListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToCapsuleListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list capsules accessible to you.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToCapsuleListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return CapsulePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Delete implements Capsule delete request.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
