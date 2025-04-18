package groups

import (
	"context"
	"net/url"
	"strings"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToGroupListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// DomainID filters the response by a domain ID.
	DomainID string `q:"domain_id"`

	// Name filters the response by group name.
	Name string `q:"name"`

	// Filters filters the response by custom filters such as
	// 'name__contains=foo'
	Filters map[string]string `q:"-"`
}

// ToGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToGroupListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}

	params := q.Query()
	for k, v := range opts.Filters {
		i := strings.Index(k, "__")
		if i > 0 && i < len(k)-2 {
			params.Add(k, v)
		} else {
			return "", InvalidListFilter{FilterName: k}
		}
	}

	q = &url.URL{RawQuery: params.Encode()}
	return q.String(), err
}

// List enumerates the Groups to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToGroupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return GroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details on a single group, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToGroupCreateMap() (map[string]any, error)
}

// CreateOpts provides options used to create a group.
type CreateOpts struct {
	// Name is the name of the new group.
	Name string `json:"name" required:"true"`

	// Description is a description of the group.
	Description string `json:"description,omitempty"`

	// DomainID is the ID of the domain the group belongs to.
	DomainID string `json:"domain_id,omitempty"`

	// Extra is free-form extra key/value pairs to describe the group.
	Extra map[string]any `json:"-"`
}

// ToGroupCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToGroupCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "group")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["group"].(map[string]any); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Create creates a new Group.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToGroupUpdateMap() (map[string]any, error)
}

// UpdateOpts provides options for updating a group.
type UpdateOpts struct {
	// Name is the name of the new group.
	Name string `json:"name,omitempty"`

	// Description is a description of the group.
	Description *string `json:"description,omitempty"`

	// DomainID is the ID of the domain the group belongs to.
	DomainID string `json:"domain_id,omitempty"`

	// Extra is free-form extra key/value pairs to describe the group.
	Extra map[string]any `json:"-"`
}

// ToGroupUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToGroupUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "group")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["group"].(map[string]any); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Update updates an existing Group.
func Update(ctx context.Context, client *gophercloud.ServiceClient, groupID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, updateURL(client, groupID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a group.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, groupID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, groupID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
