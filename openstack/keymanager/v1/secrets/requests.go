package secrets

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// DateFilter represents a valid filter to use for filtering
// secrets by their date during a list.
type DateFilter string

const (
	DateFilterGT  DateFilter = "gt"
	DateFilterGTE DateFilter = "gte"
	DateFilterLT  DateFilter = "lt"
	DateFilterLTE DateFilter = "lte"
)

// DateQuery represents a date field to be used for listing secrets.
// If no filter is specified, the query will act as if "equal" is used.
type DateQuery struct {
	Date   time.Time
	Filter DateFilter
}

// SecretType represents a valid secret type.
type SecretType string

const (
	SymmetricSecret   SecretType = "symmetric"
	PublicSecret      SecretType = "public"
	PrivateSecret     SecretType = "private"
	PassphraseSecret  SecretType = "passphrase"
	CertificateSecret SecretType = "certificate"
	OpaqueSecret      SecretType = "opaque"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToSecretListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// Offset is the starting index within the total list of the secrets that
	// you would like to retrieve.
	Offset int `q:"offset"`

	// Limit is the maximum number of records to return.
	Limit int `q:"limit"`

	// Name will select all secrets with a matching name.
	Name string `q:"name"`

	// Alg will select all secrets with a matching algorithm.
	Alg string `q:"alg"`

	// Mode will select all secrets with a matching mode.
	Mode string `q:"mode"`

	// Bits will select all secrets with a matching bit length.
	Bits int `q:"bits"`

	// SecretType will select all secrets with a matching secret type.
	SecretType SecretType `q:"secret_type"`

	// ACLOnly will select all secrets with an ACL that contains the user.
	ACLOnly *bool `q:"acl_only"`

	// CreatedQuery will select all secrets with a created date matching
	// the query.
	CreatedQuery *DateQuery

	// UpdatedQuery will select all secrets with an updated date matching
	// the query.
	UpdatedQuery *DateQuery

	// ExpirationQuery will select all secrets with an expiration date
	// matching the query.
	ExpirationQuery *DateQuery

	// Sort will sort the results in the requested order.
	Sort string `q:"sort"`
}

// ToSecretListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToSecretListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	params := q.Query()

	if opts.CreatedQuery != nil {
		created := opts.CreatedQuery.Date.Format(time.RFC3339)
		if v := opts.CreatedQuery.Filter; v != "" {
			created = fmt.Sprintf("%s:%s", v, created)
		}

		params.Add("created", created)
	}

	if opts.UpdatedQuery != nil {
		updated := opts.UpdatedQuery.Date.Format(time.RFC3339)
		if v := opts.UpdatedQuery.Filter; v != "" {
			updated = fmt.Sprintf("%s:%s", v, updated)
		}

		params.Add("updated", updated)
	}

	if opts.ExpirationQuery != nil {
		expiration := opts.ExpirationQuery.Date.Format(time.RFC3339)
		if v := opts.ExpirationQuery.Filter; v != "" {
			expiration = fmt.Sprintf("%s:%s", v, expiration)
		}

		params.Add("expiration", expiration)
	}

	q = &url.URL{RawQuery: params.Encode()}

	return q.String(), err
}

// List retrieves a list of Secrets.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToSecretListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SecretPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details of a secrets.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetPayloadOpts represents options used for obtaining a payload.
type GetPayloadOpts struct {
	PayloadContentType string `h:"Accept"`
}

// GetPayloadOptsBuilder allows extensions to add additional parameters to
// the GetPayload request.
type GetPayloadOptsBuilder interface {
	ToSecretPayloadGetParams() (map[string]string, error)
}

// ToSecretPayloadGetParams formats a GetPayloadOpts into a query string.
func (opts GetPayloadOpts) ToSecretPayloadGetParams() (map[string]string, error) {
	return gophercloud.BuildHeaders(opts)
}

// GetPayload retrieves the payload of a secret.
func GetPayload(ctx context.Context, client *gophercloud.ServiceClient, id string, opts GetPayloadOptsBuilder) (r PayloadResult) {
	h := map[string]string{"Accept": "text/plain"}

	if opts != nil {
		headers, err := opts.ToSecretPayloadGetParams()
		if err != nil {
			r.Err = err
			return
		}
		for k, v := range headers {
			h[k] = v
		}
	}

	url := payloadURL(client, id)
	resp, err := client.Get(ctx, url, nil, &gophercloud.RequestOpts{
		MoreHeaders:      h,
		OkCodes:          []int{200},
		KeepResponseBody: true,
	})
	r.Body, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToSecretCreateMap() (map[string]any, error)
}

// CreateOpts provides options used to create a secrets.
type CreateOpts struct {
	// Algorithm is the algorithm of the secret.
	Algorithm string `json:"algorithm,omitempty"`

	// BitLength is the bit length of the secret.
	BitLength int `json:"bit_length,omitempty"`

	// Mode is the mode of encryption for the secret.
	Mode string `json:"mode,omitempty"`

	// Name is the name of the secret
	Name string `json:"name,omitempty"`

	// Payload is the secret.
	Payload string `json:"payload,omitempty"`

	// PayloadContentType is the content type of the payload.
	PayloadContentType string `json:"payload_content_type,omitempty"`

	// PayloadContentEncoding is the content encoding of the payload.
	PayloadContentEncoding string `json:"payload_content_encoding,omitempty"`

	// SecretType is the type of secret.
	SecretType SecretType `json:"secret_type,omitempty"`

	// Expiration is the expiration date of the secret.
	Expiration *time.Time `json:"-"`
}

// ToSecretCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToSecretCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.Expiration != nil {
		b["expiration"] = opts.Expiration.Format(gophercloud.RFC3339NoZ)
	}

	return b, nil
}

// Create creates a new secrets.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSecretCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a secrets.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToSecretUpdateRequest() (string, map[string]string, error)
}

// UpdateOpts represents parameters to add a payload to an existing
// secret which does not already contain a payload.
type UpdateOpts struct {
	// ContentType represents the content type of the payload.
	ContentType string `h:"Content-Type"`

	// ContentEncoding represents the content encoding of the payload.
	ContentEncoding string `h:"Content-Encoding"`

	// Payload is the payload of the secret.
	Payload string
}

// ToUpdateCreateRequest formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToSecretUpdateRequest() (string, map[string]string, error) {
	h, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return "", nil, err
	}

	return opts.Payload, h, nil
}

// Update modifies the attributes of a secrets.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	url := updateURL(client, id)
	h := make(map[string]string)
	var b string

	if opts != nil {
		payload, headers, err := opts.ToSecretUpdateRequest()
		if err != nil {
			r.Err = err
			return
		}

		for k, v := range headers {
			h[k] = v
		}

		b = payload
	}

	resp, err := client.Put(ctx, url, strings.NewReader(b), nil, &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetMetadata will list metadata for a given secret.
func GetMetadata(ctx context.Context, client *gophercloud.ServiceClient, secretID string) (r MetadataResult) {
	resp, err := client.Get(ctx, metadataURL(client, secretID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// MetadataOpts is a map that contains key-value pairs for secret metadata.
type MetadataOpts map[string]string

// CreateMetadataOptsBuilder allows extensions to add additional parameters to
// the CreateMetadata request.
type CreateMetadataOptsBuilder interface {
	ToMetadataCreateMap() (map[string]any, error)
}

// ToMetadataCreateMap converts a MetadataOpts into a request body.
func (opts MetadataOpts) ToMetadataCreateMap() (map[string]any, error) {
	return map[string]any{"metadata": opts}, nil
}

// CreateMetadata will set metadata for a given secret.
func CreateMetadata(ctx context.Context, client *gophercloud.ServiceClient, secretID string, opts CreateMetadataOptsBuilder) (r MetadataCreateResult) {
	b, err := opts.ToMetadataCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, metadataURL(client, secretID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetMetadatum will get a single key/value metadata from a secret.
func GetMetadatum(ctx context.Context, client *gophercloud.ServiceClient, secretID string, key string) (r MetadatumResult) {
	resp, err := client.Get(ctx, metadatumURL(client, secretID, key), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// MetadatumOpts represents a single metadata.
type MetadatumOpts struct {
	Key   string `json:"key" required:"true"`
	Value string `json:"value" required:"true"`
}

// CreateMetadatumOptsBuilder allows extensions to add additional parameters to
// the CreateMetadatum request.
type CreateMetadatumOptsBuilder interface {
	ToMetadatumCreateMap() (map[string]any, error)
}

// ToMetadatumCreateMap converts a MetadatumOpts into a request body.
func (opts MetadatumOpts) ToMetadatumCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// CreateMetadatum will add a single key/value metadata to a secret.
func CreateMetadatum(ctx context.Context, client *gophercloud.ServiceClient, secretID string, opts CreateMetadatumOptsBuilder) (r MetadatumCreateResult) {
	b, err := opts.ToMetadatumCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, metadataURL(client, secretID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateMetadatumOptsBuilder allows extensions to add additional parameters to
// the UpdateMetadatum request.
type UpdateMetadatumOptsBuilder interface {
	ToMetadatumUpdateMap() (map[string]any, string, error)
}

// ToMetadatumUpdateMap converts a MetadataOpts into a request body.
func (opts MetadatumOpts) ToMetadatumUpdateMap() (map[string]any, string, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, opts.Key, err
}

// UpdateMetadatum will update a single key/value metadata to a secret.
func UpdateMetadatum(ctx context.Context, client *gophercloud.ServiceClient, secretID string, opts UpdateMetadatumOptsBuilder) (r MetadatumResult) {
	b, key, err := opts.ToMetadatumUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, metadatumURL(client, secretID, key), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteMetadatum will delete an individual metadatum from a secret.
func DeleteMetadatum(ctx context.Context, client *gophercloud.ServiceClient, secretID string, key string) (r MetadatumDeleteResult) {
	resp, err := client.Delete(ctx, metadatumURL(client, secretID, key), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
