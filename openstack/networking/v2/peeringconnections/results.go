package peeringconnections

import (
	"encoding/json"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*PeeringConnection, error) {
	var s PeeringConnection
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "peering_connection")
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}

type PeeringConnection struct {
	ID             string `json:"id"`
	PeerStatus     string `json:"peering_status"`
	Description    string `json:"description"`
	ConnectionType string `json:"connection_type"`
	Status         string `json:"status"`
	VpcId          string `json:"src_vpc_id"`
	PeerOrgId      string `json:"dest_org_id"`
	PeerVpcId      string `json:"dest_vpc_id"`
}

func (r *PeeringConnection) UnmarshalJSON(b []byte) error {
	type tmp PeeringConnection
	var s struct {
		tmp
	}

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	*r = PeeringConnection(s.tmp)

	return nil
}

type PeeringConnectionPage struct {
	pagination.LinkedPageBase
}

func (r PeeringConnectionPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"peering_connections_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

func (r PeeringConnectionPage) IsEmpty() (bool, error) {
	peeringConnections, err := ExtractPeeringConnections(r)
	if err != nil {
		return true, err
	}
	return len(peeringConnections) == 0, nil
}

func ExtractPeeringConnections(r pagination.Page) ([]PeeringConnection, error) {
	var s []PeeringConnection
	err := ExtractPeeringConnectionsInto(r, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func ExtractPeeringConnectionsInto(r pagination.Page, v any) error {
	return r.(PeeringConnectionPage).ExtractIntoSlicePtr(v, "peering_connections")
}
