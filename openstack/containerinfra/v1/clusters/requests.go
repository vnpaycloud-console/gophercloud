package clusters

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
	ClusterTemplateID string            `json:"cluster_template_id" required:"true"`
	CreateTimeout     *int              `json:"create_timeout"`
	DiscoveryURL      string            `json:"discovery_url,omitempty"`
	DockerVolumeSize  *int              `json:"docker_volume_size,omitempty"`
	FlavorID          string            `json:"flavor_id,omitempty"`
	Keypair           string            `json:"keypair,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	MasterCount       *int              `json:"master_count,omitempty"`
	MasterFlavorID    string            `json:"master_flavor_id,omitempty"`
	Name              string            `json:"name"`
	NodeCount         *int              `json:"node_count,omitempty"`
	FloatingIPEnabled *bool             `json:"floating_ip_enabled,omitempty"`
	MasterLBEnabled   *bool             `json:"master_lb_enabled,omitempty"`
	FixedNetwork      string            `json:"fixed_network,omitempty"`
	FixedSubnet       string            `json:"fixed_subnet,omitempty"`
	MergeLabels       *bool             `json:"merge_labels,omitempty"`
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
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a specific clusters based on its unique ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
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
	ToClustersListQuery() (string, error)
}

// ListOpts allows the sorting of paginated collections through
// the API. SortKey allows you to sort by a particular cluster attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for pagination.
type ListOpts struct {
	Marker  string `q:"marker"`
	Limit   int    `q:"limit"`
	SortKey string `q:"sort_key"`
	SortDir string `q:"sort_dir"`
}

// ToClustersListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClustersListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// clusters. It accepts a ListOptsBuilder, which allows you to sort
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToClustersListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListDetail returns a Pager which allows you to iterate over a collection of
// clusters with detailed information.
// It accepts a ListOptsBuilder, which allows you to sort the returned
// collection for greater efficiency.
func ListDetail(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listDetailURL(c)
	if opts != nil {
		query, err := opts.ToClustersListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
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
	ToClustersUpdateMap() (map[string]any, error)
}

// ToClusterUpdateMap assembles a request body based on the contents of
// UpdateOpts.
func (opts UpdateOpts) ToClustersUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Update implements cluster updated request.
func Update[T UpdateOptsBuilder](ctx context.Context, client *gophercloud.ServiceClient, id string, opts []T) (r UpdateResult) {
	var o []map[string]any
	for _, opt := range opts {
		b, err := opt.ToClustersUpdateMap()
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

type UpgradeOpts struct {
	ClusterTemplate string `json:"cluster_template" required:"true"`
	MaxBatchSize    *int   `json:"max_batch_size,omitempty"`
	NodeGroup       string `json:"nodegroup,omitempty"`
}

// UpgradeOptsBuilder allows extensions to add additional parameters to the
// Upgrade request.
type UpgradeOptsBuilder interface {
	ToClustersUpgradeMap() (map[string]any, error)
}

// ToClustersUpgradeMap constructs a request body from UpgradeOpts.
func (opts UpgradeOpts) ToClustersUpgradeMap() (map[string]any, error) {
	if opts.MaxBatchSize == nil {
		defaultMaxBatchSize := 1
		opts.MaxBatchSize = &defaultMaxBatchSize
	}
	return gophercloud.BuildRequestBody(opts, "")
}

// Upgrade implements cluster upgrade request.
func Upgrade(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpgradeOptsBuilder) (r UpgradeResult) {
	b, err := opts.ToClustersUpgradeMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, upgradeURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ResizeOptsBuilder allows extensions to add additional parameters to the
// Resize request.
type ResizeOptsBuilder interface {
	ToClusterResizeMap() (map[string]any, error)
}

// ResizeOpts params
type ResizeOpts struct {
	NodeCount     *int     `json:"node_count" required:"true"`
	NodesToRemove []string `json:"nodes_to_remove,omitempty"`
	NodeGroup     string   `json:"nodegroup,omitempty"`
}

// ToClusterResizeMap constructs a request body from ResizeOpts.
func (opts ResizeOpts) ToClusterResizeMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Resize an existing cluster node count.
func Resize(ctx context.Context, client *gophercloud.ServiceClient, id string, opts ResizeOptsBuilder) (r ResizeResult) {
	b, err := opts.ToClusterResizeMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, resizeURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
