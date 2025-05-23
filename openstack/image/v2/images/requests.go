package images

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToImageListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
//
// http://developer.openstack.org/api-ref-image-v2.html
type ListOpts struct {
	// ID is the ID of the image.
	// Multiple IDs can be specified by constructing a string
	// such as "in:uuid1,uuid2,uuid3".
	ID string `q:"id"`

	// Integer value for the limit of values to return.
	Limit int `q:"limit"`

	// UUID of the server at which you want to set a marker.
	Marker string `q:"marker"`

	// Name filters on the name of the image.
	// Multiple names can be specified by constructing a string
	// such as "in:name1,name2,name3".
	Name string `q:"name"`

	// Visibility filters on the visibility of the image.
	Visibility ImageVisibility `q:"visibility"`

	// Hidden filters on the hidden status of the image.
	Hidden bool `q:"os_hidden"`

	// MemberStatus filters on the member status of the image.
	MemberStatus ImageMemberStatus `q:"member_status"`

	// Owner filters on the project ID of the image.
	Owner string `q:"owner"`

	// Status filters on the status of the image.
	// Multiple statuses can be specified by constructing a string
	// such as "in:saving,queued".
	Status ImageStatus `q:"status"`

	// SizeMin filters on the size_min image property.
	SizeMin int64 `q:"size_min"`

	// SizeMax filters on the size_max image property.
	SizeMax int64 `q:"size_max"`

	// Sort sorts the results using the new style of sorting. See the OpenStack
	// Image API reference for the exact syntax.
	//
	// Sort cannot be used with the classic sort options (sort_key and sort_dir).
	Sort string `q:"sort"`

	// SortKey will sort the results based on a specified image property.
	SortKey string `q:"sort_key"`

	// SortDir will sort the list results either ascending or decending.
	SortDir string `q:"sort_dir"`

	// Tags filters on specific image tags.
	Tags []string `q:"tag"`

	// CreatedAtQuery filters images based on their creation date.
	CreatedAtQuery *ImageDateQuery

	// UpdatedAtQuery filters images based on their updated date.
	UpdatedAtQuery *ImageDateQuery

	// ContainerFormat filters images based on the container_format.
	// Multiple container formats can be specified by constructing a
	// string such as "in:bare,ami".
	ContainerFormat string `q:"container_format"`

	// DiskFormat filters images based on the disk_format.
	// Multiple disk formats can be specified by constructing a string
	// such as "in:qcow2,iso".
	DiskFormat string `q:"disk_format"`
}

// ToImageListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToImageListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	params := q.Query()

	if opts.CreatedAtQuery != nil {
		createdAt := opts.CreatedAtQuery.Date.Format(time.RFC3339)
		if v := opts.CreatedAtQuery.Filter; v != "" {
			createdAt = fmt.Sprintf("%s:%s", v, createdAt)
		}

		params.Add("created_at", createdAt)
	}

	if opts.UpdatedAtQuery != nil {
		updatedAt := opts.UpdatedAtQuery.Date.Format(time.RFC3339)
		if v := opts.UpdatedAtQuery.Filter; v != "" {
			updatedAt = fmt.Sprintf("%s:%s", v, updatedAt)
		}

		params.Add("updated_at", updatedAt)
	}

	q = &url.URL{RawQuery: params.Encode()}

	return q.String(), err
}

// List implements image list request.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToImageListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		imagePage := ImagePage{
			serviceURL:     c.ServiceURL(),
			LinkedPageBase: pagination.LinkedPageBase{PageResult: r},
		}

		return imagePage
	})
}

// CreateOptsBuilder allows extensions to add parameters to the Create request.
type CreateOptsBuilder interface {
	// Returns value that can be passed to json.Marshal
	ToImageCreateMap() (map[string]any, error)
}

// CreateOpts represents options used to create an image.
type CreateOpts struct {
	// Name is the name of the new image.
	Name string `json:"name" required:"true"`

	// Id is the the image ID.
	ID string `json:"id,omitempty"`

	// Visibility defines who can see/use the image.
	Visibility *ImageVisibility `json:"visibility,omitempty"`

	// Hidden is whether the image is listed in default image list or not.
	Hidden *bool `json:"os_hidden,omitempty"`

	// Tags is a set of image tags.
	Tags []string `json:"tags,omitempty"`

	// ContainerFormat is the format of the
	// container. Valid values are ami, ari, aki, bare, and ovf.
	ContainerFormat string `json:"container_format,omitempty"`

	// DiskFormat is the format of the disk. If set,
	// valid values are ami, ari, aki, vhd, vmdk, raw, qcow2, vdi,
	// and iso.
	DiskFormat string `json:"disk_format,omitempty"`

	// MinDisk is the amount of disk space in
	// GB that is required to boot the image.
	MinDisk int `json:"min_disk,omitempty"`

	// MinRAM is the amount of RAM in MB that
	// is required to boot the image.
	MinRAM int `json:"min_ram,omitempty"`

	// protected is whether the image is not deletable.
	Protected *bool `json:"protected,omitempty"`

	// properties is a set of properties, if any, that
	// are associated with the image.
	Properties map[string]string `json:"-"`
}

// ToImageCreateMap assembles a request body based on the contents of
// a CreateOpts.
func (opts CreateOpts) ToImageCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.Properties != nil {
		for k, v := range opts.Properties {
			b[k] = v
		}
	}
	return b, nil
}

// Create implements create image request.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToImageCreateMap()
	if err != nil {
		r.Err = err
		return r
	}
	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{201}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete implements image delete request.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get implements image get request.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Update implements image updated request.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToImageUpdateMap()
	if err != nil {
		r.Err = err
		return r
	}
	resp, err := client.Patch(ctx, updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/openstack-images-v2.1-json-patch"},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	// returns value implementing json.Marshaler which when marshaled matches
	// the patch schema:
	// http://specs.openstack.org/openstack/glance-specs/specs/api/v2/http-patch-image-api-v2.html
	ToImageUpdateMap() ([]any, error)
}

// UpdateOpts implements UpdateOpts
type UpdateOpts []Patch

// ToImageUpdateMap assembles a request body based on the contents of
// UpdateOpts.
func (opts UpdateOpts) ToImageUpdateMap() ([]any, error) {
	m := make([]any, len(opts))
	for i, patch := range opts {
		patchJSON := patch.ToImagePatchMap()
		m[i] = patchJSON
	}
	return m, nil
}

// Patch represents a single update to an existing image. Multiple updates
// to an image can be submitted at the same time.
type Patch interface {
	ToImagePatchMap() map[string]any
}

// UpdateVisibility represents an updated visibility property request.
type UpdateVisibility struct {
	Visibility ImageVisibility
}

// ToImagePatchMap assembles a request body based on UpdateVisibility.
func (r UpdateVisibility) ToImagePatchMap() map[string]any {
	return map[string]any{
		"op":    "replace",
		"path":  "/visibility",
		"value": r.Visibility,
	}
}

// ReplaceImageHidden represents an updated os_hidden property request.
type ReplaceImageHidden struct {
	NewHidden bool
}

// ToImagePatchMap assembles a request body based on ReplaceImageHidden.
func (r ReplaceImageHidden) ToImagePatchMap() map[string]any {
	return map[string]any{
		"op":    "replace",
		"path":  "/os_hidden",
		"value": r.NewHidden,
	}
}

// ReplaceImageName represents an updated image_name property request.
type ReplaceImageName struct {
	NewName string
}

// ToImagePatchMap assembles a request body based on ReplaceImageName.
func (r ReplaceImageName) ToImagePatchMap() map[string]any {
	return map[string]any{
		"op":    "replace",
		"path":  "/name",
		"value": r.NewName,
	}
}

// ReplaceImageChecksum represents an updated checksum property request.
type ReplaceImageChecksum struct {
	Checksum string
}

// ReplaceImageChecksum assembles a request body based on ReplaceImageChecksum.
func (r ReplaceImageChecksum) ToImagePatchMap() map[string]any {
	return map[string]any{
		"op":    "replace",
		"path":  "/checksum",
		"value": r.Checksum,
	}
}

// ReplaceImageTags represents an updated tags property request.
type ReplaceImageTags struct {
	NewTags []string
}

// ToImagePatchMap assembles a request body based on ReplaceImageTags.
func (r ReplaceImageTags) ToImagePatchMap() map[string]any {
	return map[string]any{
		"op":    "replace",
		"path":  "/tags",
		"value": r.NewTags,
	}
}

// ReplaceImageMinDisk represents an updated min_disk property request.
type ReplaceImageMinDisk struct {
	NewMinDisk int
}

// ToImagePatchMap assembles a request body based on ReplaceImageTags.
func (r ReplaceImageMinDisk) ToImagePatchMap() map[string]any {
	return map[string]any{
		"op":    "replace",
		"path":  "/min_disk",
		"value": r.NewMinDisk,
	}
}

// ReplaceImageMinRam represents an updated min_ram property request.
type ReplaceImageMinRam struct {
	NewMinRam int
}

// ToImagePatchMap assembles a request body based on ReplaceImageTags.
func (r ReplaceImageMinRam) ToImagePatchMap() map[string]any {
	return map[string]any{
		"op":    "replace",
		"path":  "/min_ram",
		"value": r.NewMinRam,
	}
}

// ReplaceImageProtected represents an updated protected property request.
type ReplaceImageProtected struct {
	NewProtected bool
}

// ToImagePatchMap assembles a request body based on ReplaceImageProtected
func (r ReplaceImageProtected) ToImagePatchMap() map[string]any {
	return map[string]any{
		"op":    "replace",
		"path":  "/protected",
		"value": r.NewProtected,
	}
}

// UpdateOp represents a valid update operation.
type UpdateOp string

const (
	AddOp     UpdateOp = "add"
	ReplaceOp UpdateOp = "replace"
	RemoveOp  UpdateOp = "remove"
)

// UpdateImageProperty represents an update property request.
type UpdateImageProperty struct {
	Op    UpdateOp
	Name  string
	Value string
}

// ToImagePatchMap assembles a request body based on UpdateImageProperty.
func (r UpdateImageProperty) ToImagePatchMap() map[string]any {
	updateMap := map[string]any{
		"op":   r.Op,
		"path": fmt.Sprintf("/%s", r.Name),
	}

	if r.Op != RemoveOp {
		updateMap["value"] = r.Value
	}

	return updateMap
}
