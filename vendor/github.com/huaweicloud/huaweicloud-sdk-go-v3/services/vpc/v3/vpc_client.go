package v3

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
)

type VpcClient struct {
	HcClient *http_client.HcHttpClient
}

func NewVpcClient(hcClient *http_client.HcHttpClient) *VpcClient {
	return &VpcClient{HcClient: hcClient}
}

func VpcClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// AddSecurityGroups 端口插入安全组
//
// 端口插入安全组
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) AddSecurityGroups(request *model.AddSecurityGroupsRequest) (*model.AddSecurityGroupsResponse, error) {
	requestDef := GenReqDefForAddSecurityGroups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddSecurityGroupsResponse), nil
	}
}

// AddSecurityGroupsInvoker 端口插入安全组
func (c *VpcClient) AddSecurityGroupsInvoker(request *model.AddSecurityGroupsRequest) *AddSecurityGroupsInvoker {
	requestDef := GenReqDefForAddSecurityGroups()
	return &AddSecurityGroupsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddSourcesToTrafficMirrorSession 流量镜像会话添加镜像源
//
// 流量镜像会话添加镜像源
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) AddSourcesToTrafficMirrorSession(request *model.AddSourcesToTrafficMirrorSessionRequest) (*model.AddSourcesToTrafficMirrorSessionResponse, error) {
	requestDef := GenReqDefForAddSourcesToTrafficMirrorSession()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddSourcesToTrafficMirrorSessionResponse), nil
	}
}

// AddSourcesToTrafficMirrorSessionInvoker 流量镜像会话添加镜像源
func (c *VpcClient) AddSourcesToTrafficMirrorSessionInvoker(request *model.AddSourcesToTrafficMirrorSessionRequest) *AddSourcesToTrafficMirrorSessionInvoker {
	requestDef := GenReqDefForAddSourcesToTrafficMirrorSession()
	return &AddSourcesToTrafficMirrorSessionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchCreateSecurityGroupRules 批量创建安全组规则
//
// 在特定安全组下批量创建安全组规则
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) BatchCreateSecurityGroupRules(request *model.BatchCreateSecurityGroupRulesRequest) (*model.BatchCreateSecurityGroupRulesResponse, error) {
	requestDef := GenReqDefForBatchCreateSecurityGroupRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchCreateSecurityGroupRulesResponse), nil
	}
}

// BatchCreateSecurityGroupRulesInvoker 批量创建安全组规则
func (c *VpcClient) BatchCreateSecurityGroupRulesInvoker(request *model.BatchCreateSecurityGroupRulesRequest) *BatchCreateSecurityGroupRulesInvoker {
	requestDef := GenReqDefForBatchCreateSecurityGroupRules()
	return &BatchCreateSecurityGroupRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchCreateSubNetworkInterface 批量创建辅助弹性网卡
//
// 批量创建辅助弹性网卡
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) BatchCreateSubNetworkInterface(request *model.BatchCreateSubNetworkInterfaceRequest) (*model.BatchCreateSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForBatchCreateSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchCreateSubNetworkInterfaceResponse), nil
	}
}

// BatchCreateSubNetworkInterfaceInvoker 批量创建辅助弹性网卡
func (c *VpcClient) BatchCreateSubNetworkInterfaceInvoker(request *model.BatchCreateSubNetworkInterfaceRequest) *BatchCreateSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForBatchCreateSubNetworkInterface()
	return &BatchCreateSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateSecurityGroup 创建安全组
//
// 创建安全组
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) CreateSecurityGroup(request *model.CreateSecurityGroupRequest) (*model.CreateSecurityGroupResponse, error) {
	requestDef := GenReqDefForCreateSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSecurityGroupResponse), nil
	}
}

// CreateSecurityGroupInvoker 创建安全组
func (c *VpcClient) CreateSecurityGroupInvoker(request *model.CreateSecurityGroupRequest) *CreateSecurityGroupInvoker {
	requestDef := GenReqDefForCreateSecurityGroup()
	return &CreateSecurityGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateSecurityGroupRule 创建安全组规则
//
// 创建安全组规则
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) CreateSecurityGroupRule(request *model.CreateSecurityGroupRuleRequest) (*model.CreateSecurityGroupRuleResponse, error) {
	requestDef := GenReqDefForCreateSecurityGroupRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSecurityGroupRuleResponse), nil
	}
}

// CreateSecurityGroupRuleInvoker 创建安全组规则
func (c *VpcClient) CreateSecurityGroupRuleInvoker(request *model.CreateSecurityGroupRuleRequest) *CreateSecurityGroupRuleInvoker {
	requestDef := GenReqDefForCreateSecurityGroupRule()
	return &CreateSecurityGroupRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateSubNetworkInterface 创建辅助弹性网卡
//
// 创建辅助弹性网卡
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) CreateSubNetworkInterface(request *model.CreateSubNetworkInterfaceRequest) (*model.CreateSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForCreateSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSubNetworkInterfaceResponse), nil
	}
}

// CreateSubNetworkInterfaceInvoker 创建辅助弹性网卡
func (c *VpcClient) CreateSubNetworkInterfaceInvoker(request *model.CreateSubNetworkInterfaceRequest) *CreateSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForCreateSubNetworkInterface()
	return &CreateSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTrafficMirrorFilter 创建流量镜像筛选条件
//
// 创建流量镜像筛选条件
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) CreateTrafficMirrorFilter(request *model.CreateTrafficMirrorFilterRequest) (*model.CreateTrafficMirrorFilterResponse, error) {
	requestDef := GenReqDefForCreateTrafficMirrorFilter()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTrafficMirrorFilterResponse), nil
	}
}

// CreateTrafficMirrorFilterInvoker 创建流量镜像筛选条件
func (c *VpcClient) CreateTrafficMirrorFilterInvoker(request *model.CreateTrafficMirrorFilterRequest) *CreateTrafficMirrorFilterInvoker {
	requestDef := GenReqDefForCreateTrafficMirrorFilter()
	return &CreateTrafficMirrorFilterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTrafficMirrorFilterRule 创建流量镜像筛选规则
//
// 创建流量镜像筛选规则
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) CreateTrafficMirrorFilterRule(request *model.CreateTrafficMirrorFilterRuleRequest) (*model.CreateTrafficMirrorFilterRuleResponse, error) {
	requestDef := GenReqDefForCreateTrafficMirrorFilterRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTrafficMirrorFilterRuleResponse), nil
	}
}

// CreateTrafficMirrorFilterRuleInvoker 创建流量镜像筛选规则
func (c *VpcClient) CreateTrafficMirrorFilterRuleInvoker(request *model.CreateTrafficMirrorFilterRuleRequest) *CreateTrafficMirrorFilterRuleInvoker {
	requestDef := GenReqDefForCreateTrafficMirrorFilterRule()
	return &CreateTrafficMirrorFilterRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTrafficMirrorSession 创建流量镜像会话
//
// 创建流量镜像会话
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) CreateTrafficMirrorSession(request *model.CreateTrafficMirrorSessionRequest) (*model.CreateTrafficMirrorSessionResponse, error) {
	requestDef := GenReqDefForCreateTrafficMirrorSession()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTrafficMirrorSessionResponse), nil
	}
}

// CreateTrafficMirrorSessionInvoker 创建流量镜像会话
func (c *VpcClient) CreateTrafficMirrorSessionInvoker(request *model.CreateTrafficMirrorSessionRequest) *CreateTrafficMirrorSessionInvoker {
	requestDef := GenReqDefForCreateTrafficMirrorSession()
	return &CreateTrafficMirrorSessionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteSecurityGroup 删除安全组
//
// 删除安全组
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) DeleteSecurityGroup(request *model.DeleteSecurityGroupRequest) (*model.DeleteSecurityGroupResponse, error) {
	requestDef := GenReqDefForDeleteSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSecurityGroupResponse), nil
	}
}

// DeleteSecurityGroupInvoker 删除安全组
func (c *VpcClient) DeleteSecurityGroupInvoker(request *model.DeleteSecurityGroupRequest) *DeleteSecurityGroupInvoker {
	requestDef := GenReqDefForDeleteSecurityGroup()
	return &DeleteSecurityGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteSecurityGroupRule 删除安全组规则
//
// 删除安全组规则
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) DeleteSecurityGroupRule(request *model.DeleteSecurityGroupRuleRequest) (*model.DeleteSecurityGroupRuleResponse, error) {
	requestDef := GenReqDefForDeleteSecurityGroupRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSecurityGroupRuleResponse), nil
	}
}

// DeleteSecurityGroupRuleInvoker 删除安全组规则
func (c *VpcClient) DeleteSecurityGroupRuleInvoker(request *model.DeleteSecurityGroupRuleRequest) *DeleteSecurityGroupRuleInvoker {
	requestDef := GenReqDefForDeleteSecurityGroupRule()
	return &DeleteSecurityGroupRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteSubNetworkInterface 删除辅助弹性网卡
//
// 删除辅助弹性网卡
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) DeleteSubNetworkInterface(request *model.DeleteSubNetworkInterfaceRequest) (*model.DeleteSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForDeleteSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSubNetworkInterfaceResponse), nil
	}
}

// DeleteSubNetworkInterfaceInvoker 删除辅助弹性网卡
func (c *VpcClient) DeleteSubNetworkInterfaceInvoker(request *model.DeleteSubNetworkInterfaceRequest) *DeleteSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForDeleteSubNetworkInterface()
	return &DeleteSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTrafficMirrorFilter 删除流量镜像筛选条件
//
// 删除流量镜像筛选条件
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) DeleteTrafficMirrorFilter(request *model.DeleteTrafficMirrorFilterRequest) (*model.DeleteTrafficMirrorFilterResponse, error) {
	requestDef := GenReqDefForDeleteTrafficMirrorFilter()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTrafficMirrorFilterResponse), nil
	}
}

// DeleteTrafficMirrorFilterInvoker 删除流量镜像筛选条件
func (c *VpcClient) DeleteTrafficMirrorFilterInvoker(request *model.DeleteTrafficMirrorFilterRequest) *DeleteTrafficMirrorFilterInvoker {
	requestDef := GenReqDefForDeleteTrafficMirrorFilter()
	return &DeleteTrafficMirrorFilterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTrafficMirrorFilterRule 删除流量镜像筛选规则
//
// 删除流量镜像筛选规则
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) DeleteTrafficMirrorFilterRule(request *model.DeleteTrafficMirrorFilterRuleRequest) (*model.DeleteTrafficMirrorFilterRuleResponse, error) {
	requestDef := GenReqDefForDeleteTrafficMirrorFilterRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTrafficMirrorFilterRuleResponse), nil
	}
}

// DeleteTrafficMirrorFilterRuleInvoker 删除流量镜像筛选规则
func (c *VpcClient) DeleteTrafficMirrorFilterRuleInvoker(request *model.DeleteTrafficMirrorFilterRuleRequest) *DeleteTrafficMirrorFilterRuleInvoker {
	requestDef := GenReqDefForDeleteTrafficMirrorFilterRule()
	return &DeleteTrafficMirrorFilterRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTrafficMirrorSession 删除流量镜像会话
//
// 删除流量镜像会话
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) DeleteTrafficMirrorSession(request *model.DeleteTrafficMirrorSessionRequest) (*model.DeleteTrafficMirrorSessionResponse, error) {
	requestDef := GenReqDefForDeleteTrafficMirrorSession()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTrafficMirrorSessionResponse), nil
	}
}

// DeleteTrafficMirrorSessionInvoker 删除流量镜像会话
func (c *VpcClient) DeleteTrafficMirrorSessionInvoker(request *model.DeleteTrafficMirrorSessionRequest) *DeleteTrafficMirrorSessionInvoker {
	requestDef := GenReqDefForDeleteTrafficMirrorSession()
	return &DeleteTrafficMirrorSessionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSecurityGroupRules 查询安全组规则列表
//
// 查询安全组规则列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ListSecurityGroupRules(request *model.ListSecurityGroupRulesRequest) (*model.ListSecurityGroupRulesResponse, error) {
	requestDef := GenReqDefForListSecurityGroupRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSecurityGroupRulesResponse), nil
	}
}

// ListSecurityGroupRulesInvoker 查询安全组规则列表
func (c *VpcClient) ListSecurityGroupRulesInvoker(request *model.ListSecurityGroupRulesRequest) *ListSecurityGroupRulesInvoker {
	requestDef := GenReqDefForListSecurityGroupRules()
	return &ListSecurityGroupRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSecurityGroups 查询安全组列表
//
// 查询某租户下的安全组列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ListSecurityGroups(request *model.ListSecurityGroupsRequest) (*model.ListSecurityGroupsResponse, error) {
	requestDef := GenReqDefForListSecurityGroups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSecurityGroupsResponse), nil
	}
}

// ListSecurityGroupsInvoker 查询安全组列表
func (c *VpcClient) ListSecurityGroupsInvoker(request *model.ListSecurityGroupsRequest) *ListSecurityGroupsInvoker {
	requestDef := GenReqDefForListSecurityGroups()
	return &ListSecurityGroupsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSubNetworkInterfaces 查询租户下辅助弹性网卡列表
//
// 查询辅助弹性网卡列表，单次查询最多返回2000条数据
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ListSubNetworkInterfaces(request *model.ListSubNetworkInterfacesRequest) (*model.ListSubNetworkInterfacesResponse, error) {
	requestDef := GenReqDefForListSubNetworkInterfaces()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSubNetworkInterfacesResponse), nil
	}
}

// ListSubNetworkInterfacesInvoker 查询租户下辅助弹性网卡列表
func (c *VpcClient) ListSubNetworkInterfacesInvoker(request *model.ListSubNetworkInterfacesRequest) *ListSubNetworkInterfacesInvoker {
	requestDef := GenReqDefForListSubNetworkInterfaces()
	return &ListSubNetworkInterfacesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTrafficMirrorFilterRules 查询流量镜像筛选规则列表
//
// 查询流量镜像筛选规则列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ListTrafficMirrorFilterRules(request *model.ListTrafficMirrorFilterRulesRequest) (*model.ListTrafficMirrorFilterRulesResponse, error) {
	requestDef := GenReqDefForListTrafficMirrorFilterRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTrafficMirrorFilterRulesResponse), nil
	}
}

// ListTrafficMirrorFilterRulesInvoker 查询流量镜像筛选规则列表
func (c *VpcClient) ListTrafficMirrorFilterRulesInvoker(request *model.ListTrafficMirrorFilterRulesRequest) *ListTrafficMirrorFilterRulesInvoker {
	requestDef := GenReqDefForListTrafficMirrorFilterRules()
	return &ListTrafficMirrorFilterRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTrafficMirrorFilters 查询流量镜像筛选条件列表
//
// 查询流量镜像筛选条件列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ListTrafficMirrorFilters(request *model.ListTrafficMirrorFiltersRequest) (*model.ListTrafficMirrorFiltersResponse, error) {
	requestDef := GenReqDefForListTrafficMirrorFilters()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTrafficMirrorFiltersResponse), nil
	}
}

// ListTrafficMirrorFiltersInvoker 查询流量镜像筛选条件列表
func (c *VpcClient) ListTrafficMirrorFiltersInvoker(request *model.ListTrafficMirrorFiltersRequest) *ListTrafficMirrorFiltersInvoker {
	requestDef := GenReqDefForListTrafficMirrorFilters()
	return &ListTrafficMirrorFiltersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTrafficMirrorSessions 查询流量镜像会话列表
//
// 查询流量镜像会话列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ListTrafficMirrorSessions(request *model.ListTrafficMirrorSessionsRequest) (*model.ListTrafficMirrorSessionsResponse, error) {
	requestDef := GenReqDefForListTrafficMirrorSessions()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTrafficMirrorSessionsResponse), nil
	}
}

// ListTrafficMirrorSessionsInvoker 查询流量镜像会话列表
func (c *VpcClient) ListTrafficMirrorSessionsInvoker(request *model.ListTrafficMirrorSessionsRequest) *ListTrafficMirrorSessionsInvoker {
	requestDef := GenReqDefForListTrafficMirrorSessions()
	return &ListTrafficMirrorSessionsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// MigrateSubNetworkInterface 迁移辅助弹性网卡
//
// 批量迁移辅助弹性网卡
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) MigrateSubNetworkInterface(request *model.MigrateSubNetworkInterfaceRequest) (*model.MigrateSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForMigrateSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.MigrateSubNetworkInterfaceResponse), nil
	}
}

// MigrateSubNetworkInterfaceInvoker 迁移辅助弹性网卡
func (c *VpcClient) MigrateSubNetworkInterfaceInvoker(request *model.MigrateSubNetworkInterfaceRequest) *MigrateSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForMigrateSubNetworkInterface()
	return &MigrateSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RemoveSecurityGroups 端口移除安全组
//
// 端口移除安全组
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) RemoveSecurityGroups(request *model.RemoveSecurityGroupsRequest) (*model.RemoveSecurityGroupsResponse, error) {
	requestDef := GenReqDefForRemoveSecurityGroups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RemoveSecurityGroupsResponse), nil
	}
}

// RemoveSecurityGroupsInvoker 端口移除安全组
func (c *VpcClient) RemoveSecurityGroupsInvoker(request *model.RemoveSecurityGroupsRequest) *RemoveSecurityGroupsInvoker {
	requestDef := GenReqDefForRemoveSecurityGroups()
	return &RemoveSecurityGroupsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RemoveSourcesFromTrafficMirrorSession 流量镜像会话移除镜像源
//
// 流量镜像会话移除镜像源
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) RemoveSourcesFromTrafficMirrorSession(request *model.RemoveSourcesFromTrafficMirrorSessionRequest) (*model.RemoveSourcesFromTrafficMirrorSessionResponse, error) {
	requestDef := GenReqDefForRemoveSourcesFromTrafficMirrorSession()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RemoveSourcesFromTrafficMirrorSessionResponse), nil
	}
}

// RemoveSourcesFromTrafficMirrorSessionInvoker 流量镜像会话移除镜像源
func (c *VpcClient) RemoveSourcesFromTrafficMirrorSessionInvoker(request *model.RemoveSourcesFromTrafficMirrorSessionRequest) *RemoveSourcesFromTrafficMirrorSessionInvoker {
	requestDef := GenReqDefForRemoveSourcesFromTrafficMirrorSession()
	return &RemoveSourcesFromTrafficMirrorSessionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowSecurityGroup 查询安全组
//
// 查询单个安全组详情
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ShowSecurityGroup(request *model.ShowSecurityGroupRequest) (*model.ShowSecurityGroupResponse, error) {
	requestDef := GenReqDefForShowSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowSecurityGroupResponse), nil
	}
}

// ShowSecurityGroupInvoker 查询安全组
func (c *VpcClient) ShowSecurityGroupInvoker(request *model.ShowSecurityGroupRequest) *ShowSecurityGroupInvoker {
	requestDef := GenReqDefForShowSecurityGroup()
	return &ShowSecurityGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowSecurityGroupRule 查询安全组规则
//
// 查询单个安全组规则
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ShowSecurityGroupRule(request *model.ShowSecurityGroupRuleRequest) (*model.ShowSecurityGroupRuleResponse, error) {
	requestDef := GenReqDefForShowSecurityGroupRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowSecurityGroupRuleResponse), nil
	}
}

// ShowSecurityGroupRuleInvoker 查询安全组规则
func (c *VpcClient) ShowSecurityGroupRuleInvoker(request *model.ShowSecurityGroupRuleRequest) *ShowSecurityGroupRuleInvoker {
	requestDef := GenReqDefForShowSecurityGroupRule()
	return &ShowSecurityGroupRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowSubNetworkInterface 查询租户下辅助弹性网卡
//
// 查询辅助弹性网卡详情
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ShowSubNetworkInterface(request *model.ShowSubNetworkInterfaceRequest) (*model.ShowSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForShowSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowSubNetworkInterfaceResponse), nil
	}
}

// ShowSubNetworkInterfaceInvoker 查询租户下辅助弹性网卡
func (c *VpcClient) ShowSubNetworkInterfaceInvoker(request *model.ShowSubNetworkInterfaceRequest) *ShowSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForShowSubNetworkInterface()
	return &ShowSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowSubNetworkInterfacesQuantity 查询租户下辅助弹性网卡数目
//
// 查询辅助弹性网卡数目
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ShowSubNetworkInterfacesQuantity(request *model.ShowSubNetworkInterfacesQuantityRequest) (*model.ShowSubNetworkInterfacesQuantityResponse, error) {
	requestDef := GenReqDefForShowSubNetworkInterfacesQuantity()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowSubNetworkInterfacesQuantityResponse), nil
	}
}

// ShowSubNetworkInterfacesQuantityInvoker 查询租户下辅助弹性网卡数目
func (c *VpcClient) ShowSubNetworkInterfacesQuantityInvoker(request *model.ShowSubNetworkInterfacesQuantityRequest) *ShowSubNetworkInterfacesQuantityInvoker {
	requestDef := GenReqDefForShowSubNetworkInterfacesQuantity()
	return &ShowSubNetworkInterfacesQuantityInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTrafficMirrorFilter 查询流量镜像筛选条件详情
//
// 查询流量镜像筛选条件详情
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ShowTrafficMirrorFilter(request *model.ShowTrafficMirrorFilterRequest) (*model.ShowTrafficMirrorFilterResponse, error) {
	requestDef := GenReqDefForShowTrafficMirrorFilter()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTrafficMirrorFilterResponse), nil
	}
}

// ShowTrafficMirrorFilterInvoker 查询流量镜像筛选条件详情
func (c *VpcClient) ShowTrafficMirrorFilterInvoker(request *model.ShowTrafficMirrorFilterRequest) *ShowTrafficMirrorFilterInvoker {
	requestDef := GenReqDefForShowTrafficMirrorFilter()
	return &ShowTrafficMirrorFilterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTrafficMirrorFilterRule 查询流量镜像筛选规则详情
//
// 查询流量镜像筛选规则详情
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ShowTrafficMirrorFilterRule(request *model.ShowTrafficMirrorFilterRuleRequest) (*model.ShowTrafficMirrorFilterRuleResponse, error) {
	requestDef := GenReqDefForShowTrafficMirrorFilterRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTrafficMirrorFilterRuleResponse), nil
	}
}

// ShowTrafficMirrorFilterRuleInvoker 查询流量镜像筛选规则详情
func (c *VpcClient) ShowTrafficMirrorFilterRuleInvoker(request *model.ShowTrafficMirrorFilterRuleRequest) *ShowTrafficMirrorFilterRuleInvoker {
	requestDef := GenReqDefForShowTrafficMirrorFilterRule()
	return &ShowTrafficMirrorFilterRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTrafficMirrorSession 查询流量镜像会话详情
//
// 查询流量镜像会话详情
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ShowTrafficMirrorSession(request *model.ShowTrafficMirrorSessionRequest) (*model.ShowTrafficMirrorSessionResponse, error) {
	requestDef := GenReqDefForShowTrafficMirrorSession()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTrafficMirrorSessionResponse), nil
	}
}

// ShowTrafficMirrorSessionInvoker 查询流量镜像会话详情
func (c *VpcClient) ShowTrafficMirrorSessionInvoker(request *model.ShowTrafficMirrorSessionRequest) *ShowTrafficMirrorSessionInvoker {
	requestDef := GenReqDefForShowTrafficMirrorSession()
	return &ShowTrafficMirrorSessionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateSecurityGroup 更新安全组
//
// 更新安全组
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) UpdateSecurityGroup(request *model.UpdateSecurityGroupRequest) (*model.UpdateSecurityGroupResponse, error) {
	requestDef := GenReqDefForUpdateSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateSecurityGroupResponse), nil
	}
}

// UpdateSecurityGroupInvoker 更新安全组
func (c *VpcClient) UpdateSecurityGroupInvoker(request *model.UpdateSecurityGroupRequest) *UpdateSecurityGroupInvoker {
	requestDef := GenReqDefForUpdateSecurityGroup()
	return &UpdateSecurityGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateSubNetworkInterface 更新辅助弹性网卡
//
// 更新辅助弹性网卡
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) UpdateSubNetworkInterface(request *model.UpdateSubNetworkInterfaceRequest) (*model.UpdateSubNetworkInterfaceResponse, error) {
	requestDef := GenReqDefForUpdateSubNetworkInterface()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateSubNetworkInterfaceResponse), nil
	}
}

// UpdateSubNetworkInterfaceInvoker 更新辅助弹性网卡
func (c *VpcClient) UpdateSubNetworkInterfaceInvoker(request *model.UpdateSubNetworkInterfaceRequest) *UpdateSubNetworkInterfaceInvoker {
	requestDef := GenReqDefForUpdateSubNetworkInterface()
	return &UpdateSubNetworkInterfaceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTrafficMirrorFilter 更新流量镜像筛选条件
//
// 更新流量镜像筛选条件
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) UpdateTrafficMirrorFilter(request *model.UpdateTrafficMirrorFilterRequest) (*model.UpdateTrafficMirrorFilterResponse, error) {
	requestDef := GenReqDefForUpdateTrafficMirrorFilter()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTrafficMirrorFilterResponse), nil
	}
}

// UpdateTrafficMirrorFilterInvoker 更新流量镜像筛选条件
func (c *VpcClient) UpdateTrafficMirrorFilterInvoker(request *model.UpdateTrafficMirrorFilterRequest) *UpdateTrafficMirrorFilterInvoker {
	requestDef := GenReqDefForUpdateTrafficMirrorFilter()
	return &UpdateTrafficMirrorFilterInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTrafficMirrorFilterRule 更新流量镜像筛选规则
//
// 更新流量镜像筛选规则
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) UpdateTrafficMirrorFilterRule(request *model.UpdateTrafficMirrorFilterRuleRequest) (*model.UpdateTrafficMirrorFilterRuleResponse, error) {
	requestDef := GenReqDefForUpdateTrafficMirrorFilterRule()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTrafficMirrorFilterRuleResponse), nil
	}
}

// UpdateTrafficMirrorFilterRuleInvoker 更新流量镜像筛选规则
func (c *VpcClient) UpdateTrafficMirrorFilterRuleInvoker(request *model.UpdateTrafficMirrorFilterRuleRequest) *UpdateTrafficMirrorFilterRuleInvoker {
	requestDef := GenReqDefForUpdateTrafficMirrorFilterRule()
	return &UpdateTrafficMirrorFilterRuleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTrafficMirrorSession 更新流量镜像会话
//
// 更新流量镜像会话
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) UpdateTrafficMirrorSession(request *model.UpdateTrafficMirrorSessionRequest) (*model.UpdateTrafficMirrorSessionResponse, error) {
	requestDef := GenReqDefForUpdateTrafficMirrorSession()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTrafficMirrorSessionResponse), nil
	}
}

// UpdateTrafficMirrorSessionInvoker 更新流量镜像会话
func (c *VpcClient) UpdateTrafficMirrorSessionInvoker(request *model.UpdateTrafficMirrorSessionRequest) *UpdateTrafficMirrorSessionInvoker {
	requestDef := GenReqDefForUpdateTrafficMirrorSession()
	return &UpdateTrafficMirrorSessionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddFirewallRules 网络ACL插入规则
//
// 网络ACL插入规则
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) AddFirewallRules(request *model.AddFirewallRulesRequest) (*model.AddFirewallRulesResponse, error) {
	requestDef := GenReqDefForAddFirewallRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddFirewallRulesResponse), nil
	}
}

// AddFirewallRulesInvoker 网络ACL插入规则
func (c *VpcClient) AddFirewallRulesInvoker(request *model.AddFirewallRulesRequest) *AddFirewallRulesInvoker {
	requestDef := GenReqDefForAddFirewallRules()
	return &AddFirewallRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AssociateSubnetFirewall 网络ACL绑定子网
//
// 网络ACL绑定子网
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) AssociateSubnetFirewall(request *model.AssociateSubnetFirewallRequest) (*model.AssociateSubnetFirewallResponse, error) {
	requestDef := GenReqDefForAssociateSubnetFirewall()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AssociateSubnetFirewallResponse), nil
	}
}

// AssociateSubnetFirewallInvoker 网络ACL绑定子网
func (c *VpcClient) AssociateSubnetFirewallInvoker(request *model.AssociateSubnetFirewallRequest) *AssociateSubnetFirewallInvoker {
	requestDef := GenReqDefForAssociateSubnetFirewall()
	return &AssociateSubnetFirewallInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateFirewall 创建网络ACL
//
// 创建网络ACL
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) CreateFirewall(request *model.CreateFirewallRequest) (*model.CreateFirewallResponse, error) {
	requestDef := GenReqDefForCreateFirewall()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateFirewallResponse), nil
	}
}

// CreateFirewallInvoker 创建网络ACL
func (c *VpcClient) CreateFirewallInvoker(request *model.CreateFirewallRequest) *CreateFirewallInvoker {
	requestDef := GenReqDefForCreateFirewall()
	return &CreateFirewallInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteFirewall 删除网络ACL
//
// 删除网络ACL
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) DeleteFirewall(request *model.DeleteFirewallRequest) (*model.DeleteFirewallResponse, error) {
	requestDef := GenReqDefForDeleteFirewall()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteFirewallResponse), nil
	}
}

// DeleteFirewallInvoker 删除网络ACL
func (c *VpcClient) DeleteFirewallInvoker(request *model.DeleteFirewallRequest) *DeleteFirewallInvoker {
	requestDef := GenReqDefForDeleteFirewall()
	return &DeleteFirewallInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DisassociateSubnetFirewall 网络ACL解绑子网
//
// 网络ACL解绑子网
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) DisassociateSubnetFirewall(request *model.DisassociateSubnetFirewallRequest) (*model.DisassociateSubnetFirewallResponse, error) {
	requestDef := GenReqDefForDisassociateSubnetFirewall()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DisassociateSubnetFirewallResponse), nil
	}
}

// DisassociateSubnetFirewallInvoker 网络ACL解绑子网
func (c *VpcClient) DisassociateSubnetFirewallInvoker(request *model.DisassociateSubnetFirewallRequest) *DisassociateSubnetFirewallInvoker {
	requestDef := GenReqDefForDisassociateSubnetFirewall()
	return &DisassociateSubnetFirewallInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListFirewall 查询网络ACL列表
//
// 查询网络ACL列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ListFirewall(request *model.ListFirewallRequest) (*model.ListFirewallResponse, error) {
	requestDef := GenReqDefForListFirewall()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListFirewallResponse), nil
	}
}

// ListFirewallInvoker 查询网络ACL列表
func (c *VpcClient) ListFirewallInvoker(request *model.ListFirewallRequest) *ListFirewallInvoker {
	requestDef := GenReqDefForListFirewall()
	return &ListFirewallInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RemoveFirewallRules 网络ACL移除规则
//
// 网络ACL移除规则
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) RemoveFirewallRules(request *model.RemoveFirewallRulesRequest) (*model.RemoveFirewallRulesResponse, error) {
	requestDef := GenReqDefForRemoveFirewallRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RemoveFirewallRulesResponse), nil
	}
}

// RemoveFirewallRulesInvoker 网络ACL移除规则
func (c *VpcClient) RemoveFirewallRulesInvoker(request *model.RemoveFirewallRulesRequest) *RemoveFirewallRulesInvoker {
	requestDef := GenReqDefForRemoveFirewallRules()
	return &RemoveFirewallRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowFirewall 查询网络ACL详情
//
// 查询网络ACL详情
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ShowFirewall(request *model.ShowFirewallRequest) (*model.ShowFirewallResponse, error) {
	requestDef := GenReqDefForShowFirewall()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowFirewallResponse), nil
	}
}

// ShowFirewallInvoker 查询网络ACL详情
func (c *VpcClient) ShowFirewallInvoker(request *model.ShowFirewallRequest) *ShowFirewallInvoker {
	requestDef := GenReqDefForShowFirewall()
	return &ShowFirewallInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateFirewall 更新网络ACL
//
// 更新网络ACL
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) UpdateFirewall(request *model.UpdateFirewallRequest) (*model.UpdateFirewallResponse, error) {
	requestDef := GenReqDefForUpdateFirewall()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateFirewallResponse), nil
	}
}

// UpdateFirewallInvoker 更新网络ACL
func (c *VpcClient) UpdateFirewallInvoker(request *model.UpdateFirewallRequest) *UpdateFirewallInvoker {
	requestDef := GenReqDefForUpdateFirewall()
	return &UpdateFirewallInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateFirewallRules 网络ACL更新规则
//
// 网络ACL更新规则
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) UpdateFirewallRules(request *model.UpdateFirewallRulesRequest) (*model.UpdateFirewallRulesResponse, error) {
	requestDef := GenReqDefForUpdateFirewallRules()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateFirewallRulesResponse), nil
	}
}

// UpdateFirewallRulesInvoker 网络ACL更新规则
func (c *VpcClient) UpdateFirewallRulesInvoker(request *model.UpdateFirewallRulesRequest) *UpdateFirewallRulesInvoker {
	requestDef := GenReqDefForUpdateFirewallRules()
	return &UpdateFirewallRulesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateAddressGroup 创建地址组
//
// 创建地址组
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) CreateAddressGroup(request *model.CreateAddressGroupRequest) (*model.CreateAddressGroupResponse, error) {
	requestDef := GenReqDefForCreateAddressGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateAddressGroupResponse), nil
	}
}

// CreateAddressGroupInvoker 创建地址组
func (c *VpcClient) CreateAddressGroupInvoker(request *model.CreateAddressGroupRequest) *CreateAddressGroupInvoker {
	requestDef := GenReqDefForCreateAddressGroup()
	return &CreateAddressGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteAddressGroup 删除地址组
//
// 删除地址组，非强制删除，删除前请确保未被其他资源引用
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) DeleteAddressGroup(request *model.DeleteAddressGroupRequest) (*model.DeleteAddressGroupResponse, error) {
	requestDef := GenReqDefForDeleteAddressGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteAddressGroupResponse), nil
	}
}

// DeleteAddressGroupInvoker 删除地址组
func (c *VpcClient) DeleteAddressGroupInvoker(request *model.DeleteAddressGroupRequest) *DeleteAddressGroupInvoker {
	requestDef := GenReqDefForDeleteAddressGroup()
	return &DeleteAddressGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteIpAddressGroupForce 强制删除地址组
//
// 强制删除地址组，删除的地址组与安全组规则关联时，会删除地址组与关联的安全组规则。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) DeleteIpAddressGroupForce(request *model.DeleteIpAddressGroupForceRequest) (*model.DeleteIpAddressGroupForceResponse, error) {
	requestDef := GenReqDefForDeleteIpAddressGroupForce()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteIpAddressGroupForceResponse), nil
	}
}

// DeleteIpAddressGroupForceInvoker 强制删除地址组
func (c *VpcClient) DeleteIpAddressGroupForceInvoker(request *model.DeleteIpAddressGroupForceRequest) *DeleteIpAddressGroupForceInvoker {
	requestDef := GenReqDefForDeleteIpAddressGroupForce()
	return &DeleteIpAddressGroupForceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAddressGroup 查询地址组列表
//
// 查询地址组列表，根据过滤条件进行过滤。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ListAddressGroup(request *model.ListAddressGroupRequest) (*model.ListAddressGroupResponse, error) {
	requestDef := GenReqDefForListAddressGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAddressGroupResponse), nil
	}
}

// ListAddressGroupInvoker 查询地址组列表
func (c *VpcClient) ListAddressGroupInvoker(request *model.ListAddressGroupRequest) *ListAddressGroupInvoker {
	requestDef := GenReqDefForListAddressGroup()
	return &ListAddressGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAddressGroup 查询地址组
//
// 查询地址组详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ShowAddressGroup(request *model.ShowAddressGroupRequest) (*model.ShowAddressGroupResponse, error) {
	requestDef := GenReqDefForShowAddressGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAddressGroupResponse), nil
	}
}

// ShowAddressGroupInvoker 查询地址组
func (c *VpcClient) ShowAddressGroupInvoker(request *model.ShowAddressGroupRequest) *ShowAddressGroupInvoker {
	requestDef := GenReqDefForShowAddressGroup()
	return &ShowAddressGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAddressGroup 更新地址组
//
// 更新地址组。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) UpdateAddressGroup(request *model.UpdateAddressGroupRequest) (*model.UpdateAddressGroupResponse, error) {
	requestDef := GenReqDefForUpdateAddressGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAddressGroupResponse), nil
	}
}

// UpdateAddressGroupInvoker 更新地址组
func (c *VpcClient) UpdateAddressGroupInvoker(request *model.UpdateAddressGroupRequest) *UpdateAddressGroupInvoker {
	requestDef := GenReqDefForUpdateAddressGroup()
	return &UpdateAddressGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AddVpcExtendCidr 添加VPC扩展网段
//
// 添加VPC的扩展网段
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) AddVpcExtendCidr(request *model.AddVpcExtendCidrRequest) (*model.AddVpcExtendCidrResponse, error) {
	requestDef := GenReqDefForAddVpcExtendCidr()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddVpcExtendCidrResponse), nil
	}
}

// AddVpcExtendCidrInvoker 添加VPC扩展网段
func (c *VpcClient) AddVpcExtendCidrInvoker(request *model.AddVpcExtendCidrRequest) *AddVpcExtendCidrInvoker {
	requestDef := GenReqDefForAddVpcExtendCidr()
	return &AddVpcExtendCidrInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateVpc 创建VPC
//
// 创建虚拟私有云
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) CreateVpc(request *model.CreateVpcRequest) (*model.CreateVpcResponse, error) {
	requestDef := GenReqDefForCreateVpc()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateVpcResponse), nil
	}
}

// CreateVpcInvoker 创建VPC
func (c *VpcClient) CreateVpcInvoker(request *model.CreateVpcRequest) *CreateVpcInvoker {
	requestDef := GenReqDefForCreateVpc()
	return &CreateVpcInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteVpc 删除VPC
//
// 删除VPC
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) DeleteVpc(request *model.DeleteVpcRequest) (*model.DeleteVpcResponse, error) {
	requestDef := GenReqDefForDeleteVpc()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteVpcResponse), nil
	}
}

// DeleteVpcInvoker 删除VPC
func (c *VpcClient) DeleteVpcInvoker(request *model.DeleteVpcRequest) *DeleteVpcInvoker {
	requestDef := GenReqDefForDeleteVpc()
	return &DeleteVpcInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListVpcs 查询VPC列表
//
// 查询vpc列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ListVpcs(request *model.ListVpcsRequest) (*model.ListVpcsResponse, error) {
	requestDef := GenReqDefForListVpcs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListVpcsResponse), nil
	}
}

// ListVpcsInvoker 查询VPC列表
func (c *VpcClient) ListVpcsInvoker(request *model.ListVpcsRequest) *ListVpcsInvoker {
	requestDef := GenReqDefForListVpcs()
	return &ListVpcsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RemoveVpcExtendCidr 移除VPC扩展网段
//
// 移除VPC扩展网段
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) RemoveVpcExtendCidr(request *model.RemoveVpcExtendCidrRequest) (*model.RemoveVpcExtendCidrResponse, error) {
	requestDef := GenReqDefForRemoveVpcExtendCidr()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RemoveVpcExtendCidrResponse), nil
	}
}

// RemoveVpcExtendCidrInvoker 移除VPC扩展网段
func (c *VpcClient) RemoveVpcExtendCidrInvoker(request *model.RemoveVpcExtendCidrRequest) *RemoveVpcExtendCidrInvoker {
	requestDef := GenReqDefForRemoveVpcExtendCidr()
	return &RemoveVpcExtendCidrInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowVpc 查询VPC详情
//
// 查询vpc详情
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) ShowVpc(request *model.ShowVpcRequest) (*model.ShowVpcResponse, error) {
	requestDef := GenReqDefForShowVpc()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowVpcResponse), nil
	}
}

// ShowVpcInvoker 查询VPC详情
func (c *VpcClient) ShowVpcInvoker(request *model.ShowVpcRequest) *ShowVpcInvoker {
	requestDef := GenReqDefForShowVpc()
	return &ShowVpcInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateVpc 更新VPC
//
// 更新vpc
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *VpcClient) UpdateVpc(request *model.UpdateVpcRequest) (*model.UpdateVpcResponse, error) {
	requestDef := GenReqDefForUpdateVpc()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateVpcResponse), nil
	}
}

// UpdateVpcInvoker 更新VPC
func (c *VpcClient) UpdateVpcInvoker(request *model.UpdateVpcRequest) *UpdateVpcInvoker {
	requestDef := GenReqDefForUpdateVpc()
	return &UpdateVpcInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
