module github.com/sbercloud-terraform/terraform-provider-sbercloud

go 1.12

require (
	github.com/hashicorp/terraform-plugin-sdk v1.14.0
	github.com/huaweicloud/golangsdk v0.0.0-20200919091224-7337da385ad9
	github.com/huaweicloud/terraform-provider-huaweicloud v1.19.1-0.20200921065647-935294670b86
)

replace github.com/huaweicloud/terraform-provider-huaweicloud => github.com/velp/terraform-provider-huaweicloud v1.19.1-0.20200921114323-443cd8d81a07
