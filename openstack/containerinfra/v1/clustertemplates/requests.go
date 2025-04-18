package clustertemplates

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]any, error)
}

// CreateOpts params
type CreateOpts struct {
	APIServerPort       *int              `json:"apiserver_port,omitempty"`
	COE                 string            `json:"coe" required:"true"`
	DNSNameServer       string            `json:"dns_nameserver,omitempty"`
	DockerStorageDriver string            `json:"docker_storage_driver,omitempty"`
	DockerVolumeSize    *int              `json:"docker_volume_size,omitempty"`
	ExternalNetworkID   string            `json:"external_network_id,omitempty"`
	FixedNetwork        string            `json:"fixed_network,omitempty"`
	FixedSubnet         string            `json:"fixed_subnet,omitempty"`
	FlavorID            string            `json:"flavor_id,omitempty"`
	FloatingIPEnabled   *bool             `json:"floating_ip_enabled,omitempty"`
	HTTPProxy           string            `json:"http_proxy,omitempty"`
	HTTPSProxy          string            `json:"https_proxy,omitempty"`
	ImageID             string            `json:"image_id" required:"true"`
	InsecureRegistry    string            `json:"insecure_registry,omitempty"`
	KeyPairID           string            `json:"keypair_id,omitempty"`
	Labels              map[string]string `json:"labels,omitempty"`
	MasterFlavorID      string            `json:"master_flavor_id,omitempty"`
	MasterLBEnabled     *bool             `json:"master_lb_enabled,omitempty"`
	Name                string            `json:"name,omitempty"`
	NetworkDriver       string            `json:"network_driver,omitempty"`
	NoProxy             string            `json:"no_proxy,omitempty"`
	Public              *bool             `json:"public,omitempty"`
	RegistryEnabled     *bool             `json:"registry_enabled,omitempty"`
	ServerType          string            `json:"server_type,omitempty"`
	TLSDisabled         *bool             `json:"tls_disabled,omitempty"`
	VolumeDriver        string            `json:"volume_driver,omitempty"`
	Hidden              *bool             `json:"hidden,omitempty"`
}

// ToClusterCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToClusterCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create requests the creation of a new cluster.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes the specified cluster ID.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToClusterTemplateListQuery() (string, error)
}

// ListOpts allows the sorting of paginated collections through
// the API. SortKey allows you to sort by a particular cluster templates attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for pagination.
type ListOpts struct {
	Marker  string `q:"marker"`
	Limit   int    `q:"limit"`
	SortKey string `q:"sort_key"`
	SortDir string `q:"sort_dir"`
}

// ToClusterTemplateListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterTemplateListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// cluster-templates. It accepts a ListOptsBuilder, which allows you to sort
// the returned collection for greater efficiency.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToClusterTemplateListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ClusterTemplatePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific cluster-template based on its unique ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateOp string

const (
	AddOp     UpdateOp = "add"
	RemoveOp  UpdateOp = "remove"
	ReplaceOp UpdateOp = "replace"
)

type UpdateOpts struct {
	Op    UpdateOp `json:"op" required:"true"`
	Path  string   `json:"path" required:"true"`
	Value any      `json:"value,omitempty"`
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToClusterTemplateUpdateMap() (map[string]any, error)
}

// ToClusterUpdateMap assembles a request body based on the contents of
// UpdateOpts.
func (opts UpdateOpts) ToClusterTemplateUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Update implements cluster updated request.
func Update[T UpdateOptsBuilder](ctx context.Context, client *gophercloud.ServiceClient, id string, opts []T) (r UpdateResult) {
	var o []map[string]any
	for _, opt := range opts {
		b, err := opt.ToClusterTemplateUpdateMap()
		if err != nil {
			r.Err = err
			return r
		}
		o = append(o, b)
	}
	resp, err := client.Patch(ctx, updateURL(client, id), o, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
