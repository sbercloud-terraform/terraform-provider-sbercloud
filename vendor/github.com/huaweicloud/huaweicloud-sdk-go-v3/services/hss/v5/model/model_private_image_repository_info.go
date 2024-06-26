package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PrivateImageRepositoryInfo repository info
type PrivateImageRepositoryInfo struct {

	// id
	Id *int64 `json:"id,omitempty"`

	// 命名空间
	Namespace *string `json:"namespace,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 镜像id
	ImageId *string `json:"image_id,omitempty"`

	// 镜像digest
	ImageDigest *string `json:"image_digest,omitempty"`

	// 镜像版本
	ImageVersion *string `json:"image_version,omitempty"`

	// 镜像类型，包含如下2种。   - private_image ：私有镜像。   - shared_image ：共享镜像。
	ImageType *string `json:"image_type,omitempty"`

	// 是否是最新版本
	LatestVersion *bool `json:"latest_version,omitempty"`

	// 扫描状态，包含如下2种。   - unscan ：未扫描。   - success ：扫描完成。   - scanning ：正在扫描。   - failed ：扫描失败。   - download_failed ：下载失败。   - image_oversized ：镜像超大。   - waiting_for_scan ：等待扫描。
	ScanStatus *string `json:"scan_status,omitempty"`

	// 镜像大小
	ImageSize *int64 `json:"image_size,omitempty"`

	// 镜像版本最后更新时间
	LatestUpdateTime *int64 `json:"latest_update_time,omitempty"`

	// 最近扫描时间
	LatestScanTime *int64 `json:"latest_scan_time,omitempty"`

	// 漏洞个数
	VulNum *int32 `json:"vul_num,omitempty"`

	// 基线扫描未通过数
	UnsafeSettingNum *int32 `json:"unsafe_setting_num,omitempty"`

	// 恶意文件数
	MaliciousFileNum *int32 `json:"malicious_file_num,omitempty"`

	// 拥有者（共享镜像参数）
	DomainName *string `json:"domain_name,omitempty"`

	// 共享镜像状态，包含如下2种。   - expired ：已过期。   - effective ：有效。
	SharedStatus *string `json:"shared_status,omitempty"`

	// 是否可扫描
	Scannable *bool `json:"scannable,omitempty"`

	// 多架构关联镜像信息
	AssociationImages *[]AssociateImages `json:"association_images,omitempty"`
}

func (o PrivateImageRepositoryInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PrivateImageRepositoryInfo struct{}"
	}

	return strings.Join([]string{"PrivateImageRepositoryInfo", string(data)}, " ")
}
