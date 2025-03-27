package capsules

import (
	"github.com/vnpaycloud-console/gophercloud/v2"
)

type ErrInvalidDataFormat struct {
	gophercloud.BaseError
}

func (e ErrInvalidDataFormat) Error() string {
	return "Data in neither json nor yaml format."
}
