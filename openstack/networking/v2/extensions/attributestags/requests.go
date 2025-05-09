package attributestags

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
)

// ReplaceAllOptsBuilder allows extensions to add additional parameters to
// the ReplaceAll request.
type ReplaceAllOptsBuilder interface {
	ToAttributeTagsReplaceAllMap() (map[string]any, error)
}

// ReplaceAllOpts provides options used to create Tags on a Resource
type ReplaceAllOpts struct {
	Tags []string `json:"tags" required:"true"`
}

// ToAttributeTagsReplaceAllMap formats a ReplaceAllOpts into the body of the
// replace request
func (opts ReplaceAllOpts) ToAttributeTagsReplaceAllMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ReplaceAll updates all tags on a resource, replacing any existing tags
func ReplaceAll(ctx context.Context, client *gophercloud.ServiceClient, resourceType string, resourceID string, opts ReplaceAllOptsBuilder) (r ReplaceAllResult) {
	b, err := opts.ToAttributeTagsReplaceAllMap()
	url := replaceURL(client, resourceType, resourceID)
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, url, &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List all tags on a resource
func List(ctx context.Context, client *gophercloud.ServiceClient, resourceType string, resourceID string) (r ListResult) {
	url := listURL(client, resourceType, resourceID)
	resp, err := client.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteAll deletes all tags on a resource
func DeleteAll(ctx context.Context, client *gophercloud.ServiceClient, resourceType string, resourceID string) (r DeleteResult) {
	url := deleteAllURL(client, resourceType, resourceID)
	resp, err := client.Delete(ctx, url, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Add a tag on a resource
func Add(ctx context.Context, client *gophercloud.ServiceClient, resourceType string, resourceID string, tag string) (r AddResult) {
	url := addURL(client, resourceType, resourceID, tag)
	resp, err := client.Put(ctx, url, nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete a tag on a resource
func Delete(ctx context.Context, client *gophercloud.ServiceClient, resourceType string, resourceID string, tag string) (r DeleteResult) {
	url := deleteURL(client, resourceType, resourceID, tag)
	resp, err := client.Delete(ctx, url, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Confirm if a tag exists on a resource
func Confirm(ctx context.Context, client *gophercloud.ServiceClient, resourceType string, resourceID string, tag string) (r ConfirmResult) {
	url := confirmURL(client, resourceType, resourceID, tag)
	resp, err := client.Get(ctx, url, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
