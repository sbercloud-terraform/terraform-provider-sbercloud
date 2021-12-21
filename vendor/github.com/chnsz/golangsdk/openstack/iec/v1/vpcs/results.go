package vpcs

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*common.VPC, error) {
	var entity common.VPC
	err := r.ExtractIntoStructPtr(&entity, "vpc")
	return &entity, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*common.VPC, error) {
	var entity common.VPC
	err := r.ExtractIntoStructPtr(&entity, "vpc")
	return &entity, err
}

type VpcPage struct {
	pagination.LinkedPageBase
}

func ExtractVpcs(r pagination.Page) ([]common.VPC, error) {
	var s struct {
		Vpcs []common.VPC `json:"vpcs"`
	}
	err := r.(VpcPage).ExtractInto(&s)
	return s.Vpcs, err
}

// IsEmpty checks whether a NetworkPage struct is empty.
func (r VpcPage) IsEmpty() (bool, error) {
	s, err := ExtractVpcs(r)
	return len(s) == 0, err
}

type UpdateResult struct {
	commonResult
}

func (r UpdateResult) Extract() (*common.VPC, error) {
	var entity common.VPC
	err := r.ExtractIntoStructPtr(&entity, "vpc")
	return &entity, err
}

type UpdateStatusResult struct {
	commonResult
}

// Vpcs vpc列表对象
type Vpcs struct {
	Vpcs  []common.VPC `json:"vpcs"`
	Count int          `json:"count"`
}

type ListResult struct {
	commonResult
}

func (r ListResult) Extract() (*Vpcs, error) {
	var entity Vpcs
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}
