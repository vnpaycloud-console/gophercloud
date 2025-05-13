package peeringconnectionapprovals

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

type ListOptsBuilder interface {
	ToPeeringConnectionApprovalListQuery() (string, error)
}

type ListOpts struct {
	PeerVPCId string `q:"src_vpc_id,omitempty"`
	PeerOrgId string `q:"src_org_id,omitempty"`
	VPCId     string `q:"dest_vpc_id,omitempty"`
	Status    string `q:"status,omitempty"`
}

func (opts ListOpts) ToPeeringConnectionApprovalListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToPeeringConnectionApprovalListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return PeeringConnectApprovalPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateOptsBuilder interface {
	ToPeeringConnectionApprovalUpdateMap() (map[string]any, error)
}

type UpdateOpts struct {
	Accept bool `json:"is_allowed,omitempty"`
}

func (opts UpdateOpts) ToPeeringConnectionApprovalUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "peering_connection_approval")
}

func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPeeringConnectionApprovalUpdateMap()
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
