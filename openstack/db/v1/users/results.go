package users

import (
	"github.com/vnpaycloud-console/gophercloud/v2"
	db "github.com/vnpaycloud-console/gophercloud/v2/openstack/db/v1/databases"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// User represents a database user
type User struct {
	// The user name
	Name string

	// The user password
	Password string

	// The databases associated with this user
	Databases []db.Database
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	gophercloud.ErrResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UserPage represents a single page of a paginated user collection.
type UserPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks to see whether the collection is empty.
func (page UserPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}

	users, err := ExtractUsers(page)
	return len(users) == 0, err
}

// NextPageURL will retrieve the next page URL.
func (page UserPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"users_links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractUsers will convert a generic pagination struct into a more
// relevant slice of User structs.
func ExtractUsers(r pagination.Page) ([]User, error) {
	var s struct {
		Users []User `json:"users"`
	}
	err := (r.(UserPage)).ExtractInto(&s)
	return s.Users, err
}
