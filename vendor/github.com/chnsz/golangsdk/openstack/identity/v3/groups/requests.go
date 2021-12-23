package groups

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToGroupListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// DomainID filters the response by a domain ID.
	DomainID string `q:"domain_id"`

	// Name filters the response by group name.
	Name string `q:"name"`
}

// ToGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToGroupListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the Groups to which the current token has access.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToGroupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return GroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// List users of the group
func ListUsers(client *golangsdk.ServiceClient, groupID string) (r UserResult) {
	_, r.Err = client.Get(listUsersURL(client, groupID), &r.Body, nil)
	return
}

// Get retrieves details on a single group, by ID.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToGroupCreateMap() (map[string]interface{}, error)
}

// CreateOpts provides options used to create a group.
type CreateOpts struct {
	// Name is the name of the new group.
	Name string `json:"name" required:"true"`

	// Description is a description of the group.
	Description string `json:"description,omitempty"`

	// DomainID is the ID of the domain the group belongs to.
	DomainID string `json:"domain_id,omitempty"`

	Extra map[string]interface{} `json:"-"`
}

// ToGroupCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToGroupCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "group")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["group"].(map[string]interface{}); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Create creates a new Group.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), &b, &r.Body, &golangsdk.RequestOpts{})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToGroupUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts provides options for updating a group.
type UpdateOpts struct {
	// Name is the name of the new group.
	Name string `json:"name,omitempty"`

	// Description is a description of the group.
	Description string `json:"description,omitempty"`

	// DomainID is the ID of the domain the group belongs to.
	DomainID string `json:"domain_id,omitempty"`

	// Extra is free-form extra key/value pairs to describe the group.
	Extra map[string]interface{} `json:"-"`
}

// ToGroupUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToGroupUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "group")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["group"].(map[string]interface{}); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Update updates an existing Group.
func Update(client *golangsdk.ServiceClient, groupID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(updateURL(client, groupID), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete deletes a group.
func Delete(client *golangsdk.ServiceClient, groupID string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, groupID), nil)
	return
}
