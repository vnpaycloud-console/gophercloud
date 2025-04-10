package users

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// List lists the existing users.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, rootURL(client), func(r pagination.PageResult) pagination.Page {
		return UserPage{pagination.SinglePageBase(r)}
	})
}

// CommonOpts are the parameters that are shared between CreateOpts and
// UpdateOpts
type CommonOpts struct {
	// Either a name or username is required. When provided, the value must be
	// unique or a 409 conflict error will be returned. If you provide a name but
	// omit a username, the latter will be set to the former; and vice versa.
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`

	// TenantID is the ID of the tenant to which you want to assign this user.
	TenantID string `json:"tenantId,omitempty"`

	// Enabled indicates whether this user is enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Email is the email address of this user.
	Email string `json:"email,omitempty"`
}

// CreateOpts represents the options needed when creating new users.
type CreateOpts CommonOpts

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToUserCreateMap() (map[string]any, error)
}

// ToUserCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToUserCreateMap() (map[string]any, error) {
	if opts.Name == "" && opts.Username == "" {
		err := gophercloud.ErrMissingInput{}
		err.Argument = "users.CreateOpts.Name/users.CreateOpts.Username"
		err.Info = "Either a Name or Username must be provided"
		return nil, err
	}
	return gophercloud.BuildRequestBody(opts, "user")
}

// Create is the operation responsible for creating new users.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToUserCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, rootURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get requests details on a single user, either by ID or Name.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, ResourceURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToUserUpdateMap() (map[string]any, error)
}

// UpdateOpts specifies the base attributes that may be updated on an
// existing server.
type UpdateOpts CommonOpts

// ToUserUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToUserUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "user")
}

// Update is the operation responsible for updating exist users by their ID.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUserUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, ResourceURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete is the operation responsible for permanently deleting a User.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, ResourceURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListRoles lists the existing roles that can be assigned to users.
func ListRoles(client *gophercloud.ServiceClient, tenantID, userID string) pagination.Pager {
	return pagination.NewPager(client, listRolesURL(client, tenantID, userID), func(r pagination.PageResult) pagination.Page {
		return RolePage{pagination.SinglePageBase(r)}
	})
}
