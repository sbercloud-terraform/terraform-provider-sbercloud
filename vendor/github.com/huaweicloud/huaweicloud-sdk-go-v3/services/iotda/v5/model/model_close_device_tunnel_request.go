package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CloseDeviceTunnelRequest Request Object
type CloseDeviceTunnelRequest struct {

	// **参数说明**：实例ID。物理多租下各实例的唯一标识，建议携带该参数，在使用专业版时必须携带该参数。您可以在IoTDA管理控制台界面，选择左侧导航栏“总览”页签查看当前实例的ID，具体获取方式请参考[[查看实例详情](https://support.huaweicloud.com/usermanual-iothub/iot_01_0121.html)](tag:hws) [[查看实例详情](https://support.huaweicloud.com/intl/zh-cn/usermanual-iothub/iot_01_0121.html)](tag:hws_hk)。
	InstanceId *string `json:"Instance-Id,omitempty"`

	// 隧道ID
	TunnelId string `json:"tunnel_id"`
}

func (o CloseDeviceTunnelRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CloseDeviceTunnelRequest struct{}"
	}

	return strings.Join([]string{"CloseDeviceTunnelRequest", string(data)}, " ")
}
