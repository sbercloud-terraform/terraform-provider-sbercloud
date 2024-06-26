package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClusterConfigurationDetailsResponse Response Object
type ShowClusterConfigurationDetailsResponse struct {

	// 指定集群配置项列表返回体，非实际返回参数
	Responses      map[string][]PackageOptions `json:"responses,omitempty"`
	HttpStatusCode int                         `json:"-"`
}

func (o ShowClusterConfigurationDetailsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterConfigurationDetailsResponse struct{}"
	}

	return strings.Join([]string{"ShowClusterConfigurationDetailsResponse", string(data)}, " ")
}
