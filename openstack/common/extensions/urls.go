package extensions

import "github.com/vnpaycloud-console/gophercloud/v2"

// ExtensionURL generates the URL for an extension resource by name.
func ExtensionURL(c *gophercloud.ServiceClient, name string) string {
	return c.ServiceURL("extensions", name)
}

// ListExtensionURL generates the URL for the extensions resource collection.
func ListExtensionURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("extensions")
}
