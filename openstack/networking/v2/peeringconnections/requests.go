package peeringconnections

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

type ListOptsBuilder interface {
	ToPeeringConnectionListQuery() (string, error)
}

type ListOpts struct {
}

func (opts ListOpts) ToPeeringConnectionListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToPeeringConnectionListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return PeeringConnectionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateOptsBuilder interface {
	ToPeeringConnectionUpdateMap() (map[string]any, error)
}

type UpdateOpts struct {
}

func (opts UpdateOpts) ToPeeringConnectionUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "peering-connection")
}

func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPeeringConnectionUpdateMap()
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

func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
