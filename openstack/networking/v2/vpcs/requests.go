package vpcs

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

type ListOptsBuilder interface {
	ToVPCListQuery() (string, error)
}

// ListOpts allows the filtering of VPCs based on their properties.
type ListOpts struct {
	Name      string `q:"name"`
	CIDR      string `q:"cidr"`
	ID        string `q:"id"`
	ProjectID string `q:"project_id"`
	Limit     int    `q:"limit"`
}

// ToVPCListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVPCListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of VPCs.
// It accepts a ListOpts struct, which allows you to filter and sort the
// returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToVPCListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return VPCPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific VPC based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToVPCCreateMap() (map[string]any, error)
}

// CreateOpts represents options used to create a VPC.
type CreateOpts struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	CIDR        string `json:"cidr,omitempty"`
}

// ToVPCCreateMap formats a CreateOpts struct into a request body.
func (opts CreateOpts) ToVPCCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "vpc")
}

// Create accepts a CreateOpts struct and creates a new VPC using the values
// provided. If the call is successful, a VPC will be returned in the
// CreateResult struct.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToVPCCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := c.Post(ctx, createURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToVPCUpdateMap() (map[string]any, error)
}

// UpdateOpts represents options used to update a VPC.
type UpdateOpts struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	EnableSNAT  *bool  `json:"enable_snat,omitempty"`
}

// ToVPCUpdateMap formats a UpdateOpts struct into a request body.
func (opts UpdateOpts) ToVPCUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "vpc")
}

// Update accepts a UpdateOpts struct and updates an existing VPC using the
// values provided. If the call is successful, a VPC will be returned in the
// UpdateResult struct.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVPCUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	h, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		r.Err = err
		return
	}

	resp, err := c.Put(ctx, updateURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{200, 201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete removes a VPC based on its unique ID. A successful response indicates
// that the VPC has been deleted.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
