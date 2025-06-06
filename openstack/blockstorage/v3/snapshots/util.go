package snapshots

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
)

// WaitForStatus will continually poll the resource, checking for a particular status.
func WaitForStatus(ctx context.Context, c *gophercloud.ServiceClient, id, status string) error {
	return gophercloud.WaitFor(ctx, func(ctx context.Context) (bool, error) {
		current, err := Get(ctx, c, id).Extract()
		if err != nil {
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}
