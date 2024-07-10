package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListImageRiskConfigRulesRequest Request Object
type ListImageRiskConfigRulesRequest struct {

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

	// 组织名称（没有镜像相关信息时，表示查询所有镜像）
	Namespace *string `json:"namespace,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 镜像版本名称
	ImageVersion *string `json:"image_version,omitempty"`

	// 基线名称
	CheckName string `json:"check_name"`

	// 标准类型，包含如下: - cn_standard : 等保合规标准 - hw_standard : 华为标准 - qt_standard : 青腾标准
	Standard string `json:"standard"`

	// 结果类型，包含如下： - pass ： 已通过 - failed : 未通过
	ResultType *string `json:"result_type,omitempty"`

	// 检查项名称，支持模糊匹配
	CheckRuleName *string `json:"check_rule_name,omitempty"`

	// 风险等级，包含如下:   - Security : 安全   - Low : 低危   - Medium : 中危   - High : 高危   - Critical : 危急
	Severity *string `json:"severity,omitempty"`
}

func (o ListImageRiskConfigRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListImageRiskConfigRulesRequest struct{}"
	}

	return strings.Join([]string{"ListImageRiskConfigRulesRequest", string(data)}, " ")
}
