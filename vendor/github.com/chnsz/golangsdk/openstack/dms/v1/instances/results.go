package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// InstanceCreate response
type InstanceCreate struct {
	InstanceID string `json:"instance_id"`
}

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

// Extract from CreateResult
func (r CreateResult) Extract() (*InstanceCreate, error) {
	var s InstanceCreate
	err := r.Result.ExtractInto(&s)
	return &s, err
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.ErrResult
}

type ListDmsResponse struct {
	Instances  []Instance `json:"instances"`
	TotalCount int        `json:"instance_num"`
}

// Instance response
type Instance struct {
	Name              string   `json:"name"`
	Engine            string   `json:"engine"`
	EngineVersion     string   `json:"engine_version"`
	Specification     string   `json:"specification"`
	PartitionNum      string   `json:"partition_num"`
	StorageSpecCode   string   `json:"storage_spec_code"`
	StorageSpace      int      `json:"storage_space"`
	UsedStorageSpace  int      `json:"used_storage_space"`
	TotalStorageSpace int      `json:"total_storage_space"`
	ConnectAddress    string   `json:"connect_address"`
	Port              int      `json:"port"`
	Nodeum            int      `json:"node_num"`
	Status            string   `json:"status"`
	SslEnable         bool     `json:"ssl_enable"`
	Description       string   `json:"description"`
	InstanceID        string   `json:"instance_id"`
	ResourceSpecCode  string   `json:"resource_spec_code"`
	Type              string   `json:"type"`
	ChargingMode      int      `json:"charging_mode"`
	VPCID             string   `json:"vpc_id"`
	VPCName           string   `json:"vpc_name"`
	CreatedAt         string   `json:"created_at"`
	ErrorCode         string   `json:"error_code"`
	ProductID         string   `json:"product_id"`
	SecurityGroupID   string   `json:"security_group_id"`
	SecurityGroupName string   `json:"security_group_name"`
	SubnetID          string   `json:"subnet_id"`
	SubnetName        string   `json:"subnet_name"`
	SubnetCIDR        string   `json:"subnet_cidr"`
	AvailableZones    []string `json:"available_zones"`
	UserID            string   `json:"user_id"`
	UserName          string   `json:"user_name"`
	OrderID           string   `json:"order_id"`
	MaintainBegin     string   `json:"maintain_begin"`
	MaintainEnd       string   `json:"maintain_end"`
}

// UpdateResult is a struct from which can get the result of update method
type UpdateResult struct {
	golangsdk.Result
}

// GetResult contains the body of getting detailed
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*Instance, error) {
	var s Instance
	err := r.Result.ExtractInto(&s)
	return &s, err
}

type DmsPage struct {
	pagination.SinglePageBase
}

func (r DmsPage) IsEmpty() (bool, error) {
	data, err := ExtractDmsInstances(r)
	if err != nil {
		return false, err
	}
	return len(data.Instances) == 0, err
}

// ExtractCloudServers is a function that takes a ListResult and returns the services' information.
func ExtractDmsInstances(r pagination.Page) (ListDmsResponse, error) {
	var s ListDmsResponse
	err := (r.(DmsPage)).ExtractInto(&s)
	return s, err
}
