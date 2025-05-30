package attachinterfaces

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// List makes a request against the nova API to list the server's interfaces.
func List(client *gophercloud.ServiceClient, serverID string) pagination.Pager {
	return pagination.NewPager(client, listInterfaceURL(client, serverID), func(r pagination.PageResult) pagination.Page {
		return InterfacePage{pagination.SinglePageBase(r)}
	})
}

// Get requests details on a single interface attachment by the server and port IDs.
func Get(ctx context.Context, client *gophercloud.ServiceClient, serverID, portID string) (r GetResult) {
	resp, err := client.Get(ctx, getInterfaceURL(client, serverID, portID), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToAttachInterfacesCreateMap() (map[string]any, error)
}

// CreateOpts specifies parameters of a new interface attachment.
type CreateOpts struct {
	// PortID is the ID of the port for which you want to create an interface.
	// The NetworkID and PortID parameters are mutually exclusive.
	// If you do not specify the PortID parameter, the OpenStack Networking API
	// v2.0 allocates a port and creates an interface for it on the network.
	PortID string `json:"port_id,omitempty"`

	// NetworkID is the ID of the network for which you want to create an interface.
	// The NetworkID and PortID parameters are mutually exclusive.
	// If you do not specify the NetworkID parameter, the OpenStack Networking
	// API v2.0 uses the network information cache that is associated with the instance.
	NetworkID string `json:"net_id,omitempty"`

	// Slice of FixedIPs. If you request a specific FixedIP address without a
	// NetworkID, the request returns a Bad Request (400) response code.
	// Note: this uses the FixedIP struct, but only the IPAddress field can be used.
	FixedIPs []FixedIP `json:"fixed_ips,omitempty"`
}

// ToAttachInterfacesCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToAttachInterfacesCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "interfaceAttachment")
}

// Create requests the creation of a new interface attachment on the server.
func Create(ctx context.Context, client *gophercloud.ServiceClient, serverID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToAttachInterfacesCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createInterfaceURL(client, serverID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete makes a request against the nova API to detach a single interface from the server.
// It needs server and port IDs to make a such request.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, serverID, portID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteInterfaceURL(client, serverID, portID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
