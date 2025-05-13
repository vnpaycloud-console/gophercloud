package peeringconnectionrequests

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

type ListOptsBuilder interface {
	ToPeeringConnectionRequestListQuery() (string, error)
}

type ListOpts struct {
}

func (opts ListOpts) ToPeeringConnectionRequestListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToPeeringConnectionRequestListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return PeeringConnectionRequestPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type CreateOptsBuilder interface {
	ToPeeringConnectionCreateMap() (map[string]any, error)
}

type CreateOpts struct {
	PeerVPCId   string `json:"dest_vpc_id,omitempty"`
	PeerOrgId   string `json:"dest_org_id,omitempty"`
	VPCId       string `json:"src_vpc_id,omitempty"`
	Description string `json:"description,omitempty"`
}

func (opts CreateOpts) ToPeeringConnectionCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "peering_connection_request")
}

func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPeeringConnectionCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := c.Post(ctx, createURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
