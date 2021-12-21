package groups

import (
	"net/http"

	"github.com/chnsz/golangsdk"
)

type CreateOpts struct {
	Name        string `json:"name,true"`
	Description string `json:"description,omitempty"`
}

type CreateOptsBuilder interface {
	ToSecurityGroupCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToSecurityGroupCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "security_group")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSecurityGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, securityGroupID string) (r DeleteResult) {
	url := DeleteURL(client, securityGroupID)
	_, r.Err = client.Delete(url, nil)
	return
}

func Get(client *golangsdk.ServiceClient, securityGroupID string) (r GetResult) {
	url := GetURL(client, securityGroupID)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

type ListOpts struct {
	Limit  int `q:"limit"`
	Offset int `q:"offset"`
}

type ListSecurityGroupsOptsBuilder interface {
	ToListSecurityGroupsQuery() (string, error)
}

func (opts ListOpts) ToListSecurityGroupsQuery() (string, error) {
	b, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func List(client *golangsdk.ServiceClient, opts ListSecurityGroupsOptsBuilder) (r ListResult) {
	listURL := rootURL(client)
	if opts != nil {
		query, err := opts.ToListSecurityGroupsQuery()
		if err != nil {
			r.Err = err
			return r
		}
		listURL += query
	}

	_, r.Err = client.Get(listURL, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}
