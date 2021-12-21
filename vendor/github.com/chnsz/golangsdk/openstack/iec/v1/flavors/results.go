package flavors

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/common"
)

type commonResult struct {
	golangsdk.Result
}

type GetResult struct {
	commonResult
}

type Flavors struct {
	Count   int             `json:"count"`
	Flavors []common.Flavor `json:"flavors"`
}

func (r GetResult) Extract() (*Flavors, error) {
	var entity Flavors
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}
