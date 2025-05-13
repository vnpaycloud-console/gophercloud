package peeringconnectionrequests

import (
	"encoding/json"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*PeeringConnectionRequest, error) {
	var s PeeringConnectionRequest
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "peering_connection_request")
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type PeeringConnectionRequest struct {
	ID             string `json:"id"`
	RequestStatus  string `json:"request_status"`
	PeerId         string `json:"peering_connection_id"`
	Description    string `json:"description"`
	ConnectionType string `json:"connection_type"`
	Status         string `json:"status"`
	VpcId          string `json:"src_vpc_id"`
	PeerOrgId      string `json:"dest_org_id"`
	PeerVpcId      string `json:"dest_vpc_id"`
}

func (r *PeeringConnectionRequest) UnmarshalJSON(b []byte) error {
	type tmp PeeringConnectionRequest
	var s struct {
		tmp
	}

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	*r = PeeringConnectionRequest(s.tmp)

	return nil
}

type PeeringConnectionRequestPage struct {
	pagination.LinkedPageBase
}

func (r PeeringConnectionRequestPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"peering_connection_requests_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

func (r PeeringConnectionRequestPage) IsEmpty() (bool, error) {
	vpcs, err := ExtractPeeringConnectionRequests(r)
	if err != nil {
		return true, err
	}
	return len(vpcs) == 0, nil
}

func ExtractPeeringConnectionRequests(r pagination.Page) ([]PeeringConnectionRequest, error) {
	var s []PeeringConnectionRequest
	err := ExtractPeeringConnectionRequestsInto(r, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func ExtractPeeringConnectionRequestsInto(r pagination.Page, v any) error {
	return r.(PeeringConnectionRequestPage).ExtractIntoSlicePtr(v, "peering_connection_requests")
}
