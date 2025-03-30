package vpcs

import (
	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a VPC resource.
func (r commonResult) Extract() (*VPC, error) {
	var s VPC
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "vpc")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a VPC.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a VPC.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a VPC.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// VPC represents, well, a VPC.
type VPC struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CIDR        string `json:"cidr"`
	SNATAddress string `json:"snat_address"`
	EnableSNAT  bool   `json:"enable_snat"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ProjectID   string `json:"project_id"`
}

func (r *VPC) UnmarshalJSON(b []byte) error {
	return nil
}

// VPCPage is the page returned by a pager when traversing over a
// collection of VPCs.
type VPCPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of VPCs has
// reached the end of a page and a new request is needed to fetch
// the next page of VPCs. It returns the URL to use for the next
// request.
func (r VPCPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"vpcs_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty returns true if a VPCPage contains no VPCs.
func (r VPCPage) IsEmpty() (bool, error) {
	vpcs, err := ExtractVPCs(r)
	if err != nil {
		return true, err
	}
	return len(vpcs) == 0, nil
}

// ExtractVPCs accepts a Page struct, specifically a VPCPage struct, and
// extracts the elements into a slice of VPC structs. In other words, a
// VPCPage contains a collection of VPC structs.
func ExtractVPCs(r pagination.Page) ([]VPC, error) {
	var s []VPC
	err := ExtractVPCsInto(r, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func ExtractVPCsInto(r pagination.Page, v any) error {
	return r.(VPCPage).ExtractIntoSlicePtr(v, "vpcs")
}
