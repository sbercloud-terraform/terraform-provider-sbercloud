package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowVulScanPolicyRequest Request Object
type ShowVulScanPolicyRequest struct {

	// 企业租户ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ShowVulScanPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowVulScanPolicyRequest struct{}"
	}

	return strings.Join([]string{"ShowVulScanPolicyRequest", string(data)}, " ")
}
