package keypairs

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/compute/v2/servers"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOptsExt adds a KeyPair option to the base CreateOpts.
type CreateOptsExt struct {
	servers.CreateOptsBuilder

	// KeyName is the name of the key pair.
	KeyName string `json:"key_name,omitempty"`
}

// ToServerCreateMap adds the key_name to the base server creation options.
func (opts CreateOptsExt) ToServerCreateMap() (map[string]interface{}, error) {
	base, err := opts.CreateOptsBuilder.ToServerCreateMap()
	if err != nil {
		return nil, err
	}

	if opts.KeyName == "" {
		return base, nil
	}

	serverMap := base["server"].(map[string]interface{})
	serverMap["key_name"] = opts.KeyName

	return base, nil
}

// List returns a Pager that allows you to iterate over a collection of KeyPairs.
func List(client *golangsdk.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return KeyPairPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToKeyPairCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies KeyPair creation or import parameters.
type CreateOpts struct {
	// Name is a friendly name to refer to this KeyPair in other services.
	Name string `json:"name" required:"true"`

	// PublicKey [optional] is a pregenerated OpenSSH-formatted public key.
	// If provided, this key will be imported and no new key will be created.
	PublicKey string `json:"public_key,omitempty"`
}

// ToKeyPairCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToKeyPairCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "keypair")
}

// Create requests the creation of a new KeyPair on the server, or to import a
// pre-existing keypair.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToKeyPairCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get returns public data about a previously uploaded KeyPair.
func Get(client *golangsdk.ServiceClient, name string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, name), &r.Body, nil)
	return
}

// Delete requests the deletion of a previous stored KeyPair from the server.
func Delete(client *golangsdk.ServiceClient, name string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, name), nil)
	return
}
