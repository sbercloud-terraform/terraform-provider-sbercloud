package tags

import (
	"github.com/chnsz/golangsdk"
)

//ResourceTags represents the struct of one autoscaling group
type ResourceTags struct {
	Tags    []ResourceTag `json:"tags"`
	SysTags []ResourceTag `json:"sys_tags"`
}

//ResourceTag is in key-value format
type ResourceTag struct {
	Key   string `json:"key" required:"ture"`
	Value string `json:"value,omitempty"`
}

//ActionResult is the action result which is the result of create or delete operations
type ActionResult struct {
	golangsdk.ErrResult
}

//GetResult contains the body of getting detailed group tags request
type GetResult struct {
	golangsdk.Result
}

//Extract method will parse the result body into ResourceTags struct
func (r GetResult) Extract() (ResourceTags, error) {
	var tags ResourceTags
	err := r.Result.ExtractInto(&tags)
	return tags, err
}

//ListResult contains the body of getting detailed group tags request
type ListResult struct {
	golangsdk.Result
}

//Extract method will parse the result body into ResourceTags struct
func (r ListResult) Extract() (ResourceTags, error) {
	var tags ResourceTags
	err := r.Result.ExtractInto(&tags)
	return tags, err
}
