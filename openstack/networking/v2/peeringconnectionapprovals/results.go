package peeringconnectionapprovals

import (
	"encoding/json"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*PeeringConnectApproval, error) {
	var s PeeringConnectApproval
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "peering_connection_approval")
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type PeeringConnectApproval struct {
	ID          string `json:"id"`
	PeerId      string `json:"peering_connection_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	VpcId       string `json:"src_vpc_id"`
	PeerOrgId   string `json:"dest_org_id"`
	PeerVpcId   string `json:"dest_vpc_id"`
	Status      string `json:"status"`
}

func (r *PeeringConnectApproval) UnmarshalJSON(b []byte) error {
	type tmp PeeringConnectApproval
	var s struct {
		tmp
	}

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	*r = PeeringConnectApproval(s.tmp)

	return nil
}

type PeeringConnectApprovalPage struct {
	pagination.LinkedPageBase
}

func (r PeeringConnectApprovalPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"peering_connection_approvals_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

func (r PeeringConnectApprovalPage) IsEmpty() (bool, error) {
	vpcs, err := ExtractPeeringConnectApprovals(r)
	if err != nil {
		return true, err
	}
	return len(vpcs) == 0, nil
}

func ExtractPeeringConnectApprovals(r pagination.Page) ([]PeeringConnectApproval, error) {
	var s []PeeringConnectApproval
	err := ExtractPeeringConnectApprovalsInto(r, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func ExtractPeeringConnectApprovalsInto(r pagination.Page, v any) error {
	return r.(PeeringConnectApprovalPage).ExtractIntoSlicePtr(v, "peering_connection_approvals")
}
