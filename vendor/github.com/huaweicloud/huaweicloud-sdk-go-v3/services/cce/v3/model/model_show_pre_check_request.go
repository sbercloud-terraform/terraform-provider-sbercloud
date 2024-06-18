package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPreCheckRequest Request Object
type ShowPreCheckRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 升级任务ID，调用集群升级API后从响应体中uid字段获取。
	TaskId string `json:"task_id"`
}

func (o ShowPreCheckRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPreCheckRequest struct{}"
	}

	return strings.Join([]string{"ShowPreCheckRequest", string(data)}, " ")
}
