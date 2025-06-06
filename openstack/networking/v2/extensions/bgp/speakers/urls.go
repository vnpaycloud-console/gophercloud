package speakers

import "github.com/vnpaycloud-console/gophercloud/v2"

const urlBase = "bgp-speakers"

// return /v2.0/bgp-speakers/{bgp-speaker-id}
func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(urlBase, id)
}

// return /v2.0/bgp-speakers
func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(urlBase)
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}
func getURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgp-speakers
func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

// return /v2.0/bgp-speakers
func createURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

// return /v2.0/bgp-speakers/{bgp-peer-id}
func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgp-speakers/{bgp-peer-id}
func updateURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}/add_bgp_peer
func addBGPPeerURL(c *gophercloud.ServiceClient, speakerID string) string {
	return c.ServiceURL(urlBase, speakerID, "add_bgp_peer")
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}/remove_bgp_peer
func removeBGPPeerURL(c *gophercloud.ServiceClient, speakerID string) string {
	return c.ServiceURL(urlBase, speakerID, "remove_bgp_peer")
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}/get_advertised_routes
func getAdvertisedRoutesURL(c *gophercloud.ServiceClient, speakerID string) string {
	return c.ServiceURL(urlBase, speakerID, "get_advertised_routes")
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}/add_gateway_network
func addGatewayNetworkURL(c *gophercloud.ServiceClient, speakerID string) string {
	return c.ServiceURL(urlBase, speakerID, "add_gateway_network")
}

// return /v2.0/bgp-speakers/{bgp-speaker-id}/remove_gateway_network
func removeGatewayNetworkURL(c *gophercloud.ServiceClient, speakerID string) string {
	return c.ServiceURL(urlBase, speakerID, "remove_gateway_network")
}
