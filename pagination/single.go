package pagination

import (
	"fmt"
	"reflect"

	"github.com/vnpaycloud-console/gophercloud/v2"
)

// SinglePageBase may be embedded in a Page that contains all of the results from an operation at once.
type SinglePageBase PageResult

// NextPageURL always returns "" to indicate that there are no more pages to return.
func (current SinglePageBase) NextPageURL() (string, error) {
	return "", nil
}

// IsEmpty satisifies the IsEmpty method of the Page interface
func (current SinglePageBase) IsEmpty() (bool, error) {
	if b, ok := current.Body.([]any); ok {
		return len(b) == 0, nil
	}
	err := gophercloud.ErrUnexpectedType{}
	err.Expected = "[]any"
	err.Actual = fmt.Sprintf("%v", reflect.TypeOf(current.Body))
	return true, err
}

// GetBody returns the single page's body. This method is needed to satisfy the
// Page interface.
func (current SinglePageBase) GetBody() any {
	return current.Body
}
