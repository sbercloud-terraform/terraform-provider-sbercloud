package stacks

import (
	"reflect"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToStackCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// The name of the stack. It must start with an alphabetic character.
	Name string `json:"stack_name" required:"true"`
	// A structure that contains either the template file or url. Call the
	// associated methods to extract the information relevant to send in a create request.
	TemplateOpts *Template `json:"-" required:"true"`
	// Enables or disables deletion of all stack resources when a stack
	// creation fails. Default is true, meaning all resources are not deleted when
	// stack creation fails.
	DisableRollback *bool `json:"disable_rollback,omitempty"`
	// A structure that contains details for the environment of the stack.
	EnvironmentOpts *Environment `json:"-"`
	// User-defined parameters to pass to the template.
	Parameters map[string]string `json:"parameters,omitempty"`
	// The timeout for stack creation in minutes.
	Timeout int `json:"timeout_mins,omitempty"`
	// A list of tags to assosciate with the Stack
	Tags []string `json:"-"`
}

// ToStackCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToStackCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if err := opts.TemplateOpts.Parse(); err != nil {
		return nil, err
	}

	if err := opts.TemplateOpts.getFileContents(opts.TemplateOpts.Parsed, ignoreIfTemplate, true); err != nil {
		return nil, err
	}
	opts.TemplateOpts.fixFileRefs()
	b["template"] = string(opts.TemplateOpts.Bin)

	files := make(map[string]string)
	for k, v := range opts.TemplateOpts.Files {
		files[k] = v
	}

	if opts.EnvironmentOpts != nil {
		if err := opts.EnvironmentOpts.Parse(); err != nil {
			return nil, err
		}
		if err := opts.EnvironmentOpts.getRRFileContents(ignoreIfEnvironment); err != nil {
			return nil, err
		}
		opts.EnvironmentOpts.fixFileRefs()
		for k, v := range opts.EnvironmentOpts.Files {
			files[k] = v
		}
		b["environment"] = string(opts.EnvironmentOpts.Bin)
	}

	if len(files) > 0 {
		b["files"] = files
	}

	if opts.Tags != nil {
		b["tags"] = strings.Join(opts.Tags, ",")
	}

	return b, nil
}

// Create accepts a CreateOpts struct and creates a new stack using the values
// provided.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToStackCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// SortDir is a type for specifying in which direction to sort a list of stacks.
type SortDir string

// SortKey is a type for specifying by which key to sort a list of stacks.
type SortKey string

var (
	// SortAsc is used to sort a list of stacks in ascending order.
	SortAsc SortDir = "asc"
	// SortDesc is used to sort a list of stacks in descending order.
	SortDesc SortDir = "desc"
	// SortName is used to sort a list of stacks by name.
	SortName SortKey = "name"
	// SortStatus is used to sort a list of stacks by status.
	SortStatus SortKey = "status"
	// SortCreatedAt is used to sort a list of stacks by date created.
	SortCreatedAt SortKey = "created_at"
	// SortUpdatedAt is used to sort a list of stacks by date updated.
	SortUpdatedAt SortKey = "updated_at"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToStackListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the rts attributes you want to see returned. SortKey allows you to sort
// by a particular network attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	ID      string  `q:"id"`
	Status  string  `q:"status"`
	Name    string  `q:"name"`
	Marker  string  `q:"marker"`
	Limit   int     `q:"limit"`
	SortKey SortKey `q:"sort_keys"`
	SortDir SortDir `q:"sort_dir"`
}

// ToStackListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToStackListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

func List(c *golangsdk.ServiceClient, opts ListOpts) ([]ListedStack, error) {
	u := listURL(c)
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return StackPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allStacks, err := ExtractStacks(pages)
	if err != nil {
		return nil, err
	}

	return FilterStacks(allStacks, opts)
}

func FilterStacks(stacks []ListedStack, opts ListOpts) ([]ListedStack, error) {

	var refinedStacks []ListedStack
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["ID"] = opts.ID
	}
	if opts.Name != "" {
		m["Name"] = opts.Name
	}
	if opts.Status != "" {
		m["Status"] = opts.Status
	}

	if len(m) > 0 && len(stacks) > 0 {
		for _, stack := range stacks {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&stack, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedStacks = append(refinedStacks, stack)
			}
		}

	} else {
		refinedStacks = stacks
	}

	return refinedStacks, nil
}

func getStructField(v *ListedStack, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

func Get(c *golangsdk.ServiceClient, stackName string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, stackName), &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the Update operation in this package.
type UpdateOptsBuilder interface {
	ToStackUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// A structure that contains either the template file or url. Call the
	// associated methods to extract the information relevant to send in a create request.
	TemplateOpts *Template `json:"-" required:"true"`
	// A structure that contains details for the environment of the stack.
	EnvironmentOpts *Environment `json:"-"`
	// User-defined parameters to pass to the template.
	Parameters map[string]string `json:"parameters,omitempty"`
	// The timeout for stack creation in minutes.
	Timeout int `json:"timeout_mins,omitempty"`
	// Enables or disables deletion of all stack resources when a stack
	// creation fails. Default is true, meaning all resources are not deleted when
	// stack creation fails.
	DisableRollback *bool `json:"disable_rollback,omitempty"`
	// A list of tags to assosciate with the Stack
	Tags []string `json:"-"`
}

// ToStackUpdateMap casts a CreateOpts struct to a map.
func (opts UpdateOpts) ToStackUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if err := opts.TemplateOpts.Parse(); err != nil {
		return nil, err
	}

	if err := opts.TemplateOpts.getFileContents(opts.TemplateOpts.Parsed, ignoreIfTemplate, true); err != nil {
		return nil, err
	}
	opts.TemplateOpts.fixFileRefs()
	b["template"] = string(opts.TemplateOpts.Bin)

	files := make(map[string]string)
	for k, v := range opts.TemplateOpts.Files {
		files[k] = v
	}

	if opts.EnvironmentOpts != nil {
		if err := opts.EnvironmentOpts.Parse(); err != nil {
			return nil, err
		}
		if err := opts.EnvironmentOpts.getRRFileContents(ignoreIfEnvironment); err != nil {
			return nil, err
		}
		opts.EnvironmentOpts.fixFileRefs()
		for k, v := range opts.EnvironmentOpts.Files {
			files[k] = v
		}
		b["environment"] = string(opts.EnvironmentOpts.Bin)
	}

	if len(files) > 0 {
		b["files"] = files
	}

	if opts.Tags != nil {
		b["tags"] = strings.Join(opts.Tags, ",")
	}

	return b, nil
}

// Update accepts an UpdateOpts struct and updates an existing stack using the values
// provided.
func Update(c *golangsdk.ServiceClient, stackName, stackID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToStackUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, stackName, stackID), b, nil, nil)
	return
}

// Delete deletes a stack based on the stack name and stack ID.
func Delete(c *golangsdk.ServiceClient, stackName, stackID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, stackName, stackID), nil)
	return
}
