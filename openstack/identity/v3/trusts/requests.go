package trusts

import (
	"context"
	"time"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToTrustCreateMap() (map[string]any, error)
}

// CreateOpts provides options used to create a new trust.
type CreateOpts struct {
	// Impersonation allows the trustee to impersonate the trustor.
	Impersonation bool `json:"impersonation"`

	// TrusteeUserID is a user who is capable of consuming the trust.
	TrusteeUserID string `json:"trustee_user_id" required:"true"`

	// TrustorUserID is a user who created the trust.
	TrustorUserID string `json:"trustor_user_id" required:"true"`

	// AllowRedelegation enables redelegation of a trust.
	AllowRedelegation bool `json:"allow_redelegation,omitempty"`

	// ExpiresAt sets expiration time on trust.
	ExpiresAt *time.Time `json:"-"`

	// ProjectID identifies the project.
	ProjectID string `json:"project_id,omitempty"`

	// RedelegationCount specifies a depth of the redelegation chain.
	RedelegationCount int `json:"redelegation_count,omitempty"`

	// RemainingUses specifies how many times a trust can be used to get a token.
	RemainingUses int `json:"remaining_uses,omitempty"`

	// Roles specifies roles that need to be granted to trustee.
	Roles []Role `json:"roles,omitempty"`
}

// ToTrustCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToTrustCreateMap() (map[string]any, error) {
	parent := "trust"
	b, err := gophercloud.BuildRequestBody(opts, parent)
	if err != nil {
		return nil, err
	}

	if opts.ExpiresAt != nil {
		if v, ok := b[parent].(map[string]any); ok {
			v["expires_at"] = opts.ExpiresAt.Format(gophercloud.RFC3339Milli)
		}
	}

	return b, nil
}

type ListOptsBuilder interface {
	ToTrustListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// TrustorUserID filters the response by a trustor user Id.
	TrustorUserID string `q:"trustor_user_id"`

	// TrusteeUserID filters the response by a trustee user Id.
	TrusteeUserID string `q:"trustee_user_id"`
}

// ToTrustListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToTrustListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Create creates a new Trust.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTrustCreateMap()
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

// Delete deletes a Trust.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, trustID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, trustID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List enumerates the Trust to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToTrustListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return TrustPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details on a single Trust, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, resourceURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListRoles lists roles delegated by a Trust.
func ListRoles(client *gophercloud.ServiceClient, id string) pagination.Pager {
	url := listRolesURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RolesPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetRole retrieves details on a single role delegated by a Trust.
func GetRole(ctx context.Context, client *gophercloud.ServiceClient, id string, roleID string) (r GetRoleResult) {
	resp, err := client.Get(ctx, getRoleURL(client, id, roleID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CheckRole checks whether a role ID is delegated by a Trust.
func CheckRole(ctx context.Context, client *gophercloud.ServiceClient, id string, roleID string) (r CheckRoleResult) {
	resp, err := client.Head(ctx, getRoleURL(client, id, roleID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
