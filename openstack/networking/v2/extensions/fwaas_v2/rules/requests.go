package rules

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

type (
	// Protocol represents a valid rule protocol
	Protocol string
)

const (
	// ProtocolAny is to allow any protocol
	ProtocolAny Protocol = "any"

	// ProtocolICMP is to allow the ICMP protocol
	ProtocolICMP Protocol = "icmp"

	// ProtocolTCP is to allow the TCP protocol
	ProtocolTCP Protocol = "tcp"

	// ProtocolUDP is to allow the UDP protocol
	ProtocolUDP Protocol = "udp"
)

type (
	// Action represents a valid rule protocol
	Action string
)

const (
	// ActionAllow is to allow traffic
	ActionAllow Action = "allow"

	// ActionDeny is to deny traffic
	ActionDeny Action = "deny"

	// ActionTCP is to reject traffic
	ActionReject Action = "reject"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToRuleListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the Firewall rule attributes you want to see returned. SortKey allows you to
// sort by a particular firewall rule attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	TenantID             string   `q:"tenant_id"`
	Name                 string   `q:"name"`
	Description          string   `q:"description"`
	Protocol             Protocol `q:"protocol"`
	Action               Action   `q:"action"`
	IPVersion            int      `q:"ip_version"`
	SourceIPAddress      string   `q:"source_ip_address"`
	DestinationIPAddress string   `q:"destination_ip_address"`
	SourcePort           string   `q:"source_port"`
	DestinationPort      string   `q:"destination_port"`
	Enabled              *bool    `q:"enabled"`
	ID                   string   `q:"id"`
	Shared               *bool    `q:"shared"`
	ProjectID            string   `q:"project_id"`
	FirewallPolicyID     string   `q:"firewall_policy_id"`
	Limit                int      `q:"limit"`
	Marker               string   `q:"marker"`
	SortKey              string   `q:"sort_key"`
	SortDir              string   `q:"sort_dir"`
}

// ToRuleListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRuleListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns a Pager which allows you to iterate over a collection of
// firewall rules. It accepts a ListOpts struct, which allows you to filter
// and sort the returned collection for greater efficiency.
//
// Default policy settings return only those firewall rules that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)

	if opts != nil {
		query, err := opts.ToRuleListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return RulePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToRuleCreateMap() (map[string]any, error)
}

// CreateOpts contains all the values needed to create a new firewall rule.
type CreateOpts struct {
	Protocol             Protocol              `json:"protocol" required:"true"`
	Action               Action                `json:"action" required:"true"`
	TenantID             string                `json:"tenant_id,omitempty"`
	ProjectID            string                `json:"project_id,omitempty"`
	Name                 string                `json:"name,omitempty"`
	Description          string                `json:"description,omitempty"`
	IPVersion            gophercloud.IPVersion `json:"ip_version,omitempty"`
	SourceIPAddress      string                `json:"source_ip_address,omitempty"`
	DestinationIPAddress string                `json:"destination_ip_address,omitempty"`
	SourcePort           string                `json:"source_port,omitempty"`
	DestinationPort      string                `json:"destination_port,omitempty"`
	Shared               *bool                 `json:"shared,omitempty"`
	Enabled              *bool                 `json:"enabled,omitempty"`
}

// ToRuleCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToRuleCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "firewall_rule")
	if err != nil {
		return nil, err
	}

	if m := b["firewall_rule"].(map[string]any); m["protocol"] == "any" {
		m["protocol"] = nil
	}

	return b, nil
}

// Create accepts a CreateOpts struct and uses the values to create a new firewall rule
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a particular firewall rule based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToRuleUpdateMap() (map[string]any, error)
}

// UpdateOpts contains the values used when updating a firewall rule.
type UpdateOpts struct {
	Protocol             *Protocol              `json:"protocol,omitempty"`
	Action               *Action                `json:"action,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	Description          *string                `json:"description,omitempty"`
	IPVersion            *gophercloud.IPVersion `json:"ip_version,omitempty"`
	SourceIPAddress      *string                `json:"source_ip_address,omitempty"`
	DestinationIPAddress *string                `json:"destination_ip_address,omitempty"`
	SourcePort           *string                `json:"source_port,omitempty"`
	DestinationPort      *string                `json:"destination_port,omitempty"`
	Shared               *bool                  `json:"shared,omitempty"`
	Enabled              *bool                  `json:"enabled,omitempty"`
}

// ToRuleUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToRuleUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "firewall_rule")
}

// Update allows firewall policies to be updated.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRuleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will permanently delete a particular firewall rule based on its unique ID.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
