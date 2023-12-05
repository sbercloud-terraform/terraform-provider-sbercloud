---
subcategory: "Identity and Access Management (IAM)"
---

# sbercloud_identity_role

Use this data source to get details of the specified IAM **system-defined** role or policy.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

The Role in Terraform is the same as Policy on console. However,
The policy name is the display name of Role, the Role name cannot
be found on Console. Please refer to the following table to configuration Role:

Display Name | Role/Policy Name | Description
---- |------------------| ---
Server Administrator | server_adm       | Server Administrator
ECS Admin | system_all_0     | All permissions of ECS service
ECS User | system_all_3     | Common permissions of ECS service, except installation, delete, reinstallation and so on.
ECS Viewer | system_all_1     | The read-only permissions to all ECS resources, which can be used for statistics and survey
IMS Administrator | ims_adm          | IMS Administrator
IMS Viewer | system_all_6     | The read-only permissions to all IMS resources, which can be used for statistics and survey.
IMS Admin| system_all_4     | All permissions of  Image Management Service
AutoScaling Administrator | as_adm           | AutoScaling Administrator
AutoScaling Admin | system_all_12    | All permissions template of AutoScaling Service.
AutoScaling FullAccess | system_all_117   | Full permissions for Auto Scaling.
AutoScaling Viewer | system_all_22    | The read-only permissions to all AutoScaling resources, which can be used for statistics and survey.
EVS Admin| system_all_7     | 	All permissions of EVS service.
EVS Viewer | system_all_2     | The read-only permissions to all EVS resorces, which can be sed for statistics and srvey.
SFS Administrator | sfs_adm          | 	SFS Administrator
SFS Admin | system_all_42    | All permissions of Scalable File Service.
SFS Viewer | system_all_43    | The read-only permissions to all Scalable File Service resources.
SFS Turbo Viewer | system_all_66    | The read-only permissions to all Scalable File Service (SFS Turbo) resources.
SFS Turbo Admin | system_all_65    | All permissions of Scalable File Service (SFS Turbo).
OBS Administrator | system_all_86    | Object Storage Service Administrator
OBS Operator | system_all_26    | Basic operation permissions to view the bucket list, obtain bucket metadata, list objects in a bucket, query bucket location, upload objects, download objects, delete objects, and obtain object ACLs
OBS Viewer| system_all_25    | Permissions to view the bucket list, obtain bucket metadata, list objects in a bucket, and query bucket location
OBS Buckets Viewer | obs_b_list       | Permissions to view the bucket list, obtain bucket metadata, and query bucket location
VPC Administrator | vpc_netadm       | VPC Administrator
VPC Admin | system_all_8     | All permissions of VPC service
VPC Viewer | system_all_5     | The read-only permissions to all VPC resources, which can be used for statistics and survey
ELB Admin | system_all_24    | All permissions of ELB service.
ELB Service Administrator | elb_adm          | ELB Service Administrator
ELB Viewer | system_all_18    | The read-only permissions to all ELB resources, which can be used for statistics and survey.
DNS Administrator | dns_adm          | DNS Administrator
DNS Admin | system_all_23    | DNS administrator permissions, which allow users to perform all operations, including creating, deleting, querying, and modifying DNS resources.
DNS Viewer | system_all_21    | Read-only permissions, which only allow users to query DNS resources.
NAT Gateway Administrator | nat_adm          | NAT Gateway Administrator
NAT Admin | system_all_13    | All permissions of NAT Gateway service
NAT Viewer | system_all_16    | The read-only permissions to all NAT Gateway resources
VPCEndpoint Administrator | vpcep_adm        | VPCEndpoint service enables you to privately connect your VPC to supported services.
RDS Administrator | rds_adm          | RDS Administrator
RDS Admin | system_all_50    | All permissions of RDS service.
RDS Viewer | system_all_48    | The read-only permissions to all RDS resources, which can be used for statistics and survey
RDS DBA | system_all_49    | DBA permissions of RDS service, except delete.
RDS FullAccess | system_all_96    | Full permissions for Relational Database Service
RDS ManageAccess | system_all_98    | Database administrator permissions for all operations except deleting RDS resources
RDS ReadOnlyAccess | system_all_97    | Read-only permissions for Relational Database Service.
DDS Administrator | dds_adm          | Document Database Service Administrator
DDS DBA | system_all_47    | DBA permissions of DDS service, except delete.
DDS FullAccess | system_all_99    | Full permissions for Document Database Service.
DDS Viewer | system_all_51    | Read-only permissions for Document Database Service.
DDS ManageAccess | system_all_100   | Database administrator permissions for all operations except deleting DDS resources.
DDS FullAccess | system_all_59    | Full permissions for Document Database Service.
DDS ReadOnlyAccess | system_all_60    | Read-only permissions for Document Database Service.
DDS Admin | system_all_52    | All permissions of DDS service.
CCE Administrator | cce_adm          | CCE Administrator
CCE FullAccess | system_all_62    | Common operation permissions on CCE cluster resources, excluding the namespace-level permissions for the clusters (with Kubernetes RBAC enabled) and the privileged administrator operations, such as agency configuration and cluster certificate generation
CCE ReadOnlyAccess | system_all_71    | Permissions to view CCE cluster resources, excluding the namespace-level permissions of the clusters (with Kubernetes RBAC enabled)
CSS Administrator| css_adm          | Cloud Search Service Administrator
CSS FullAccess | system_all_104   | All permissions for Cloud Search Service
CSS ReadOnlyAccess | system_all_103   | Read-only permissions for viewing Cloud Search Service
ServiceStage Administrator | svcstg_adm       | ServiceStage administrator, who has full permissions for this service
ServiceStage Developer | svcstg_dev       | ServiceStage developer, who has full permissions for this service but does not have the permission for creating infrastructure
ServiceStage Developer | system_all_27    | Developer permissions of ServiceStage service(exclude review and approve).
ServiceStage Operator | svcstg_opr       | ServiceStage operator, who has the read-only permission for this service
APM Administrator | apm_adm          | Application Performance Monitor Admin
CES Administrator | ces_adm          | CloudEye Service Administrator
CTS Administrator | cts_adm          | CloudTrace Service Administrator
DCS Administrator | dcs_admin        | Distributed Cache Service Administrator
DIS Administrator | dis_adm          | Data Ingestion Service User
KMS Administrator | kms_adm          | KMS Administrator
MRS Administrator | mrs_adm          | MRS Administrator
SWR Admin | swr_adm          | Software Repository Administrator
SMN Administrator | smn_adm          | SMN Administrator
TMS Administrator | tms_adm          | Tag Management Service Administrator
Security Administrator | secu_admin       | Full permissions for Identity and Access Management. This role does not have permissions for switching roles.
Tenant Administrator | te_admin         | Tenant Administrator (Exclude IAM)
Tenant Guest | readonly         | Tenant Guest (Exclude IAM)
EPS FullAccess | system_all_63    | All permissions of EPS service.

## Example Usage

```hcl
data "sbercloud_identity_role" "auth_admin" {
  name = "secu_admin"
}
```

## Argument Reference

* `display_name` - (Optional, String) Specifies the display name of the role displayed on the console.
  It is recommended to use this parameter instead of `name` and required if `name` is not specified.

* `name` - (Optional, String) Specifies the name of the role for internal use.
  It's required if `display_name` is not specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.
* `description` - The description of the policy.
* `catalog` - The service catalog of the policy.
* `type` - The display mode of the policy.
* `policy` - The content of the policy.
