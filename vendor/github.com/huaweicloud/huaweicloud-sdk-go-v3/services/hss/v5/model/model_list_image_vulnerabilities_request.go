package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListImageVulnerabilitiesRequest Request Object
type ListImageVulnerabilitiesRequest struct {

	// region id
	Region string `json:"region"`

	// 租户企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 镜像类型，包含如下:   - private_image : 私有镜像仓库   - shared_image : 共享镜像仓库
	ImageType string `json:"image_type"`

	// 偏移量：指定返回记录的开始位置，必须为数字，取值范围为大于或等于0，默认0
	Offset *int32 `json:"offset,omitempty"`

	// 每页显示个数
	Limit *int32 `json:"limit,omitempty"`

	// 镜像id
	ImageId string `json:"image_id"`

	// 企业仓库实例ID，swr共享版无需使用该参数
	InstanceId *string `json:"instance_id,omitempty"`

	// 组织名称
	Namespace string `json:"namespace"`

	// 镜像名称
	ImageName string `json:"image_name"`

	// 镜像版本名称
	TagName string `json:"tag_name"`

	// 危险程度，包含如下3种。   - immediate_repair ：高危。   - delay_repair ：中危。   - not_needed_repair ：低危。
	RepairNecessity *string `json:"repair_necessity,omitempty"`

	// 漏洞ID（支持模糊查询）
	VulId *string `json:"vul_id,omitempty"`

	// 软件名
	AppName *string `json:"app_name,omitempty"`
}

func (o ListImageVulnerabilitiesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListImageVulnerabilitiesRequest struct{}"
	}

	return strings.Join([]string{"ListImageVulnerabilitiesRequest", string(data)}, " ")
}
