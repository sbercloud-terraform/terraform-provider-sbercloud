package acceptance

import (
	"encoding/json"
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

var (
	SBC_REGION_NAME = os.Getenv("SBC_REGION_NAME")

	SBC_ENTERPRISE_PROJECT_ID      = os.Getenv("SBC_ENTERPRISE_PROJECT_ID")
	SBC_ENTERPRISE_PROJECT_ID_TEST = os.Getenv("SBC_ENTERPRISE_PROJECT_ID_TEST")
	SBC_PROJECT_ID                 = os.Getenv("SBC_PROJECT_ID")

	SBC_DEPRECATED_ENVIRONMENT = os.Getenv("SBC_DEPRECATED_ENVIRONMENT")

	SBC_ADMIN       = os.Getenv("SBC_ADMIN")
	SBC_DOMAIN_ID   = os.Getenv("SBC_DOMAIN_ID")
	SBC_DOMAIN_NAME = os.Getenv("SBC_DOMAIN_NAME")

	SBC_ACCESS_KEY = os.Getenv("SBC_ACCESS_KEY")
	SBC_SECRET_KEY = os.Getenv("SBC_SECRET_KEY")

	SBC_DLI_FLINK_JAR_OBS_PATH = os.Getenv("SBC_DLI_FLINK_JAR_OBS_PATH")

	SBC_SWR_SHARING_ACCOUNT = os.Getenv("SBC_SWR_SHARING_ACCOUNT")

	SBC_FGS_TRIGGER_LTS_AGENCY             = os.Getenv("SBC_FGS_TRIGGER_LTS_AGENCY")
	SBC_OBS_BUCKET_NAME                    = os.Getenv("SBC_OBS_BUCKET_NAME")
	SBC_DWS_MUTIL_AZS                      = os.Getenv("SBC_DWS_MUTIL_AZS")
	SBC_ENTERPRISE_MIGRATE_PROJECT_ID_TEST = os.Getenv("SBC_ENTERPRISE_MIGRATE_PROJECT_ID_TEST")

	SBC_CHARGING_MODE   = os.Getenv("SBC_CHARGING_MODE")
	SBC_KMS_ENVIRONMENT = os.Getenv("SBC_KMS_ENVIRONMENT")

	SBC_SHARED_BACKUP_ID     = os.Getenv("SBC_SHARED_BACKUP_ID")
	SBC_DEST_PROJECT_ID      = os.Getenv("SBC_DEST_PROJECT_ID")
	SBC_DEST_PROJECT_ID_TEST = os.Getenv("SBC_DEST_PROJECT_ID_TEST")
	SBC_DEST_REGION          = os.Getenv("SBC_DEST_REGION")

	SBC_CERTIFICATE_CONTENT                = os.Getenv("SBC_CERTIFICATE_CONTENT")
	SBC_CERTIFICATE_CONTENT_UPDATE         = os.Getenv("SBC_CERTIFICATE_CONTENT_UPDATE")
	SBC_CERTIFICATE_PRIVATE_KEY            = os.Getenv("SBC_CERTIFICATE_PRIVATE_KEY")
	SBC_NEW_CERTIFICATE_CONTENT            = os.Getenv("SBC_NEW_CERTIFICATE_CONTENT")
	SBC_NEW_CERTIFICATE_PRIVATE_KEY        = os.Getenv("SBC_NEW_CERTIFICATE_PRIVATE_KEY")
	SBC_CERTIFICATE_ROOT_CA                = os.Getenv("SBC_CERTIFICATE_ROOT_CA")
	SBC_NEW_CERTIFICATE_ROOT_CA            = os.Getenv("SBC_NEW_CERTIFICATE_ROOT_CA")
	SBC_GM_CERTIFICATE_CONTENT             = os.Getenv("SBC_GM_CERTIFICATE_CONTENT")
	SBC_GM_CERTIFICATE_PRIVATE_KEY         = os.Getenv("SBC_GM_CERTIFICATE_PRIVATE_KEY")
	SBC_GM_ENC_CERTIFICATE_CONTENT         = os.Getenv("SBC_GM_ENC_CERTIFICATE_CONTENT")
	SBC_GM_ENC_CERTIFICATE_PRIVATE_KEY     = os.Getenv("SBC_GM_ENC_CERTIFICATE_PRIVATE_KEY")
	SBC_GM_CERTIFICATE_CHAIN               = os.Getenv("SBC_GM_CERTIFICATE_CHAIN")
	SBC_NEW_GM_CERTIFICATE_CONTENT         = os.Getenv("SBC_NEW_GM_CERTIFICATE_CONTENT")
	SBC_NEW_GM_CERTIFICATE_PRIVATE_KEY     = os.Getenv("SBC_NEW_GM_CERTIFICATE_PRIVATE_KEY")
	SBC_NEW_GM_ENC_CERTIFICATE_CONTENT     = os.Getenv("SBC_NEW_GM_ENC_CERTIFICATE_CONTENT")
	SBC_NEW_GM_ENC_CERTIFICATE_PRIVATE_KEY = os.Getenv("SBC_NEW_GM_ENC_CERTIFICATE_PRIVATE_KEY")
	SBC_NEW_GM_CERTIFICATE_CHAIN           = os.Getenv("SBC_NEW_GM_CERTIFICATE_CHAIN")
	SBC_CODEARTS_RESOURCE_POOL_ID          = os.Getenv("SBC_CODEARTS_RESOURCE_POOL_ID")

	SBC_RDS_INSTANCE_ID = os.Getenv("SBC_RDS_INSTANCE_ID")

	SBC_OBS_ENDPOINT        = os.Getenv("SBC_OBS_ENDPOINT")
	SBC_SFS_TURBO_BACKUP_ID = os.Getenv("SBC_SFS_TURBO_BACKUP_ID")

	SBC_VPC_ID            = os.Getenv("SBC_VPC_ID")
	SBC_SUBNET_ID         = os.Getenv("SBC_SUBNET_ID")
	SBC_SECURITY_GROUP_ID = os.Getenv("SBC_SECURITY_GROUP_ID")

	SBC_DDS_SECOND_LEVEL_MONITORING_ENABLED = os.Getenv("SBC_DDS_SECOND_LEVEL_MONITORING_ENABLED")

	SBC_APIG_DEDICATED_INSTANCE_ID             = os.Getenv("SBC_APIG_DEDICATED_INSTANCE_ID")
	SBC_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID = os.Getenv("SBC_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID")
	SBC_FGS_AGENCY_NAME                        = os.Getenv("SBC_FGS_AGENCY_NAME")
)

// TestAccProviderFactories is a static map containing only the main provider instance
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

// TestAccProvider is the "main" provider instance
var TestAccProvider *schema.Provider

func init() {
	TestAccProvider = sbercloud.Provider()

	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"sbercloud": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
	}
}

// ServiceFunc the SberCloud resource query functions.
type ServiceFunc func(*config.Config, *terraform.ResourceState) (interface{}, error)

// resourceCheck resource check object, only used in the package.
type resourceCheck struct {
	resourceName    string
	resourceObject  interface{}
	getResourceFunc ServiceFunc
	resourceType    string
}

const (
	resourceTypeCode   = "resource"
	dataSourceTypeCode = "dataSource"

	checkAttrRegexpStr = `^\$\{([^\}]+)\}$`
)

/*
InitDataSourceCheck build a 'resourceCheck' object. Only used to check datasource attributes.

	Parameters:
	  resourceName:    The resource name is used to check in the terraform.State.e.g. : sbercloud_waf_domain.domain_1.
	Return:
	  *resourceCheck: resourceCheck object
*/
func InitDataSourceCheck(sourceName string) *resourceCheck {
	return &resourceCheck{
		resourceName: sourceName,
		resourceType: dataSourceTypeCode,
	}
}

/*
InitResourceCheck build a 'resourceCheck' object. The common test methods are provided in 'resourceCheck'.

	Parameters:
	  resourceName:    The resource name is used to check in the terraform.State.e.g. : sbercloud_waf_domain.domain_1.
	  resourceObject:  Resource object, used to check whether the resource exists in SberCloud.
	  getResourceFunc: The function used to get the resource object.
	Return:
	  *resourceCheck: resourceCheck object
*/
func InitResourceCheck(resourceName string, resourceObject interface{}, getResourceFunc ServiceFunc) *resourceCheck {
	return &resourceCheck{
		resourceName:    resourceName,
		resourceObject:  resourceObject,
		getResourceFunc: getResourceFunc,
		resourceType:    resourceTypeCode,
	}
}

func parseVariableToName(varStr string) (string, string, error) {
	var resName, keyName string
	// Check the format of the variable.
	match, _ := regexp.MatchString(checkAttrRegexpStr, varStr)
	if !match {
		return resName, keyName, fmtp.Errorf("The type of 'variable' is error, "+
			"expected ${resourceType.name.field} got %s", varStr)
	}

	reg, err := regexp.Compile(checkAttrRegexpStr)
	if err != nil {
		return resName, keyName, fmtp.Errorf("The acceptance function is wrong.")
	}
	mArr := reg.FindStringSubmatch(varStr)
	if len(mArr) != 2 {
		return resName, keyName, fmtp.Errorf("The type of 'variable' is error, "+
			"expected ${resourceType.name.field} got %s", varStr)
	}

	// Get resName and keyName from variable.
	strs := strings.Split(mArr[1], ".")
	for i, s := range strs {
		if strings.Contains(s, "sbercloud_") {
			resName = strings.Join(strs[0:i+2], ".")
			keyName = strings.Join(strs[i+2:], ".")
			break
		}
	}
	return resName, keyName, nil
}

/*
TestCheckResourceAttrWithVariable validates the variable in state for the given name/key combination.

	Parameters:
	  resourceName: The resource name is used to check in the terraform.State.
	  key:          The field name of the resource.
	  variable:     The variable name of the value to be checked.

	  variable such like ${sbercloud_waf_certificate.certificate_1.id}
	  or ${data.sbercloud_waf_policies.policies_2.policies.0.id}
*/
func TestCheckResourceAttrWithVariable(resourceName, key, varStr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resName, keyName, err := parseVariableToName(varStr)
		if err != nil {
			return err
		}

		if strings.EqualFold(resourceName, resName) {
			return fmtp.Errorf("Meaningless verification. " +
				"The referenced resource cannot be the current resource.")
		}

		// Get the value based on resName and keyName from the state.
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmtp.Errorf("Can't find %s in state : %s.", resName, ok)
		}
		value := rs.Primary.Attributes[keyName]

		return resource.TestCheckResourceAttr(resourceName, key, value)(s)
	}
}

// CheckResourceDestroy check whether resources destroyed in SberCloud.
func (rc *resourceCheck) CheckResourceDestroy() resource.TestCheckFunc {
	if strings.Compare(rc.resourceType, dataSourceTypeCode) == 0 {
		fmtp.Errorf("Error, you built a resourceCheck with 'InitDataSourceCheck', " +
			"it cannot run CheckResourceDestroy().")
		return nil
	}
	return func(s *terraform.State) error {
		strs := strings.Split(rc.resourceName, ".")
		var resourceType string
		for _, str := range strs {
			if strings.Contains(str, "sbercloud_") {
				resourceType = strings.Trim(str, " ")
				break
			}
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			conf := TestAccProvider.Meta().(*config.Config)
			if rc.getResourceFunc != nil {
				if _, err := rc.getResourceFunc(conf, rs); err == nil {
					return fmtp.Errorf("failed to destroy resource. The resource of %s : %s still exists.",
						resourceType, rs.Primary.ID)
				}
			} else {
				return fmtp.Errorf("The 'getResourceFunc' is nil, please set it during initialization.")
			}
		}
		return nil
	}
}

// CheckResourceExists check whether resources exist in SberCloud.
//
//	func (rc *resourceCheck) CheckResourceExists() resource.TestCheckFunc {
//		return func(s *terraform.State) error {
//			rs, ok := s.RootModule().Resources[rc.resourceName]
//			if !ok {
//				return fmtp.Errorf("Can not found the resource or data source in state: %s", rc.resourceName)
//			}
//			if rs.Primary.ID == "" {
//				return fmtp.Errorf("No id set for the resource or data source: %s", rc.resourceName)
//			}
//			if strings.EqualFold(rc.resourceType, dataSourceTypeCode) {
//				return nil
//			}
//
//			if rc.getResourceFunc != nil {
//				conf := TestAccProvider.Meta().(*config.Config)
//				r, err := rc.getResourceFunc(conf, rs)
//				if err != nil {
//					return fmtp.Errorf("checking resource %s %s exists error: %s ",
//						rc.resourceName, rs.Primary.ID, err)
//				}
//				if rc.resourceObject != nil {
//					b, err := json.Marshal(r)
//					if err != nil {
//						return fmtp.Errorf("marshaling resource %s %s error: %s ",
//							rc.resourceName, rs.Primary.ID, err)
//					}
//					json.Unmarshal(b, rc.resourceObject)
//				} else {
//					logp.Printf("[WARN] The 'resourceObject' is nil, please set it during initialization.")
//				}
//			} else {
//				return fmtp.Errorf("The 'getResourceFunc' is nil, please set it.")
//			}
//
//			return nil
//		}
//	}
func (rc *resourceCheck) checkResourceExists(s *terraform.State) error {
	rs, ok := s.RootModule().Resources[rc.resourceName]
	if !ok {
		return fmt.Errorf("can not found the resource or data source in state: %s", rc.resourceName)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No id set for the resource or data source: %s", rc.resourceName)
	}
	if strings.EqualFold(rc.resourceType, dataSourceTypeCode) {
		return nil
	}

	if rc.getResourceFunc == nil {
		return fmt.Errorf("the 'getResourceFunc' is nil, please set it during initialization")
	}

	conf := TestAccProvider.Meta().(*config.Config)
	r, err := rc.getResourceFunc(conf, rs)
	if err != nil {
		return fmt.Errorf("checking resource %s %s exists error: %s ",
			rc.resourceName, rs.Primary.ID, err)
	}

	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshaling resource %s %s error: %s ",
			rc.resourceName, rs.Primary.ID, err)
	}

	// unmarshal the response body into the resourceObject
	if rc.resourceObject != nil {
		return json.Unmarshal(b, rc.resourceObject)
	}

	return nil
}

// CheckResourceExists check whether resources exist
func (rc *resourceCheck) CheckResourceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return rc.checkResourceExists(s)
	}
}

func (rc *resourceCheck) CheckMultiResourcesExists(count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error
		for i := 0; i < count; i++ {
			rcCopy := *rc
			rcCopy.resourceName = fmt.Sprintf("%s.%d", rcCopy.resourceName, i)
			err = rcCopy.checkResourceExists(s)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func preCheckRequiredEnvVars(t *testing.T) {
	if SBC_REGION_NAME == "" {
		t.Fatal("SBC_REGION_NAME must be set for acceptance tests")
	}
}

func TestAccPreCheck(t *testing.T) {
	preCheckRequiredEnvVars(t)
}

func TestAccPreCheckVpcId(t *testing.T) {
	if SBC_VPC_ID == "" {
		t.Skip("SBC_VPC_ID must be set for the acceptance test")
	}
}

func TestAccPreCheckSubnetId(t *testing.T) {
	if SBC_SUBNET_ID == "" {
		t.Skip("SBC_SUBNET_ID must be set for the acceptance test")
	}
}

func TestAccPreCheckSecurityGroupId(t *testing.T) {
	if SBC_SECURITY_GROUP_ID == "" {
		t.Skip("SBC_SECURITY_GROUP_ID must be set for the acceptance test")
	}
}

func TestAccPreCheckRdsInstanceId(t *testing.T) {
	if SBC_RDS_INSTANCE_ID == "" {
		t.Skip("SBC_RDS_INSTANCE_ID must be set for RDS acceptance tests")
	}
}

func TestAccPreCheckOBSEndpoint(t *testing.T) {
	if SBC_OBS_ENDPOINT == "" {
		t.Skip("SBC_OBS_ENDPOINT must be set for the acceptance test")
	}
}

func TestAccPrecheckSFSTurboBackupId(t *testing.T) {
	if SBC_SFS_TURBO_BACKUP_ID == "" {
		t.Skip("SBC_SFS_TURBO_BACKUP_ID must be set for the acceptance test")
	}
}

func TestAccPreCheckAcceptBackup(t *testing.T) {
	if SBC_SHARED_BACKUP_ID == "" {
		t.Skip("SBC_SHARED_BACKUP_ID must be set for CBR backup share acceptance")
	}
}

func TestAccPreCheckDestProjectIds(t *testing.T) {
	if SBC_DEST_PROJECT_ID == "" || SBC_DEST_PROJECT_ID_TEST == "" {
		t.Skip("SBC_DEST_PROJECT_ID and SBC_DEST_PROJECT_ID_TEST must be set for acceptance test.")
	}
}

func TestAccPreCheckReplication(t *testing.T) {
	if SBC_DEST_REGION == "" || SBC_DEST_PROJECT_ID == "" {
		t.Skip("Skip the replication policy acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckChargingMode(t *testing.T) {
	if SBC_CHARGING_MODE != "prePaid" {
		t.Skip("This environment does not support prepaid tests")
	}
}

func TestAccPreCheckDeprecated(t *testing.T) {
	if SBC_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}

	preCheckRequiredEnvVars(t)
}

func TestAccPreCheckEpsID(t *testing.T) {
	if SBC_ENTERPRISE_PROJECT_ID == "" {
		t.Skip("This environment does not support Enterprise Project ID tests")
	}
}

func TestAccPreCheckProject(t *testing.T) {
	if SBC_ENTERPRISE_PROJECT_ID_TEST == "" {
		t.Skip("This environment does not support project tests")
	}
}

func TestAccPreCheckAdminOnly(t *testing.T) {
	if SBC_ADMIN == "" {
		t.Skip("Skipping test because it requires the admin privileges")
	}
}

func TestAccPreCheckOBS(t *testing.T) {
	if SBC_ACCESS_KEY == "" || SBC_SECRET_KEY == "" {
		t.Skip("SBC_ACCESS_KEY and SBC_SECRET_KEY must be set for OBS acceptance tests")
	}
}

func TestAccPreCheckSWRDomian(t *testing.T) {
	if SBC_SWR_SHARING_ACCOUNT == "" {
		t.Skip("SBC_SWR_SHARING_ACCOUNT must be set for swr domian tests, " +
			"the value of SBC_SWR_SHARING_ACCOUNT should be another IAM user name")
	}
}

func TestAccPreCheckFgsTrigger(t *testing.T) {
	if SBC_FGS_TRIGGER_LTS_AGENCY == "" {
		t.Skip("SBC_FGS_TRIGGER_LTS_AGENCY must be set for FGS trigger acceptance tests")
	}
}

func TestAccPreCheckKms(t *testing.T) {
	if SBC_KMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support KMS tests")
	}
}

func TestAccPreCheckOBSBucket(t *testing.T) {
	if SBC_OBS_BUCKET_NAME == "" {
		t.Skip("SBC_OBS_BUCKET_NAME must be set for OBS object acceptance tests")
	}
}

func TestAccPrecheckDomainId(t *testing.T) {
	if SBC_DOMAIN_ID == "" {
		t.Skip("SBC_DOMAIN_ID must be set for acceptance tests")
	}
}
func TestAccPreCheckProjectID(t *testing.T) {
	if SBC_PROJECT_ID == "" {
		t.Skip("SBC_PROJECT_ID must be set for acceptance tests")
	}
}

func TestAccPreCheckMutilAZ(t *testing.T) {
	if SBC_DWS_MUTIL_AZS == "" {
		t.Skip("SBC_DWS_MUTIL_AZS must be set for the acceptance test")
	}
}
func TestAccPreCheckMigrateEpsID(t *testing.T) {
	if SBC_ENTERPRISE_PROJECT_ID_TEST == "" || SBC_ENTERPRISE_MIGRATE_PROJECT_ID_TEST == "" {
		t.Skip("The environment variables does not support Migrate Enterprise Project ID for acc tests")
	}
}

func TestAccPreCheckUpdateCertificateContent(t *testing.T) {
	if SBC_CERTIFICATE_CONTENT == "" || SBC_CERTIFICATE_CONTENT_UPDATE == "" {
		t.Skip("SBC_CERTIFICATE_CONTENT, SBC_CERTIFICATE_CONTENT_UPDATE must be set for this test")
	}
}

// lintignore:AT003
func TestAccPreCheckCertificateWithoutRootCA(t *testing.T) {
	if SBC_CERTIFICATE_CONTENT == "" || SBC_CERTIFICATE_PRIVATE_KEY == "" ||
		SBC_NEW_CERTIFICATE_CONTENT == "" || SBC_NEW_CERTIFICATE_PRIVATE_KEY == "" {
		t.Skip("SBC_CERTIFICATE_CONTENT, SBC_CERTIFICATE_PRIVATE_KEY, SBC_NEW_CERTIFICATE_CONTENT and " +
			"SBC_NEW_CERTIFICATE_PRIVATE_KEY must be set for simple acceptance tests of SSL certificate resource")
	}
}

// lintignore:AT003
func TestAccPreCheckDDSSecondLevelMonitoringEnabled(t *testing.T) {
	if SBC_DDS_SECOND_LEVEL_MONITORING_ENABLED == "" {
		t.Skip("SBC_DDS_SECOND_LEVEL_MONITORING_ENABLED must be set for the acceptance test")
	}
}

// lintignore:AT003
func TestAccPreCheckApigSubResourcesRelatedInfo(t *testing.T) {
	if SBC_APIG_DEDICATED_INSTANCE_ID == "" {
		t.Skip("Before running APIG acceptance tests, please ensure the env 'HW_APIG_DEDICATED_INSTANCE_ID' has been configured")
	}
}

// lintignore:AT003
func TestAccPreCheckApigChannelRelatedInfo(t *testing.T) {
	if SBC_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID == "" {
		t.Skip("Before running APIG acceptance tests, please ensure the env 'HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID' has been configured")
	}
}

// lintignore:AT003
func TestAccPreCheckFgsAgency(t *testing.T) {
	// The agency should be FunctionGraph and authorize these roles:
	// For the acceptance tests of the async invoke configuration:
	// + FunctionGraph FullAccess
	// + DIS Operator
	// + OBS Administrator
	// + SMN Administrator
	// For the acceptance tests of the function trigger and the application:
	// + LTS Administrator
	if SBC_FGS_AGENCY_NAME == "" {
		t.Skip("HW_FGS_AGENCY_NAME must be set for FGS acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckCertificateFull(t *testing.T) {
	TestAccPreCheckCertificateWithoutRootCA(t)
	if SBC_CERTIFICATE_ROOT_CA == "" || SBC_NEW_CERTIFICATE_ROOT_CA == "" {
		t.Skip("SBC_CERTIFICATE_ROOT_CA and SBC_NEW_CERTIFICATE_ROOT_CA must be set for root CA validation")
	}
}

// lintignore:AT003
func TestAccPreCheckGMCertificate(t *testing.T) {
	if SBC_GM_CERTIFICATE_CONTENT == "" || SBC_GM_CERTIFICATE_PRIVATE_KEY == "" ||
		SBC_GM_ENC_CERTIFICATE_CONTENT == "" || SBC_GM_ENC_CERTIFICATE_PRIVATE_KEY == "" ||
		SBC_GM_CERTIFICATE_CHAIN == "" ||
		SBC_NEW_GM_CERTIFICATE_CONTENT == "" || SBC_NEW_GM_CERTIFICATE_PRIVATE_KEY == "" ||
		SBC_NEW_GM_ENC_CERTIFICATE_CONTENT == "" || SBC_NEW_GM_ENC_CERTIFICATE_PRIVATE_KEY == "" ||
		SBC_NEW_GM_CERTIFICATE_CHAIN == "" {
		t.Skip("SBC_GM_CERTIFICATE_CONTENT, SBC_GM_CERTIFICATE_PRIVATE_KEY, SBC_GM_ENC_CERTIFICATE_CONTENT," +
			" SBC_GM_ENC_CERTIFICATE_PRIVATE_KEY, SBC_GM_CERTIFICATE_CHAIN, SBC_NEW_GM_CERTIFICATE_CONTENT," +
			" SBC_NEW_GM_CERTIFICATE_PRIVATE_KEY, SBC_NEW_GM_ENC_CERTIFICATE_CONTENT," +
			" SBC_NEW_GM_ENC_CERTIFICATE_PRIVATE_KEY, SBC_NEW_GM_CERTIFICATE_CHAIN must be set for GM certificate")
	}
}

// lintignore:AT003
func TestAccPreCheckCodeArtsDeployResourcePoolID(t *testing.T) {
	if SBC_CODEARTS_RESOURCE_POOL_ID == "" {
		t.Skip("SBC_CODEARTS_RESOURCE_POOL_ID must be set for this acceptance test")
	}
}

func RandomAccResourceName() string {
	return fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
}

func RandomAccResourceNameWithDash() string {
	return fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
}

func RandomCidr() string {
	return fmt.Sprintf("172.16.%d.0/24", acctest.RandIntRange(0, 255))
}

func RandomCidrAndGatewayIp() (string, string) {
	seed := acctest.RandIntRange(0, 255)
	return fmt.Sprintf("172.16.%d.0/24", seed), fmt.Sprintf("172.16.%d.1", seed)
}

func RandomPassword(customChars ...string) string {
	var specialChars string
	if len(customChars) < 1 {
		specialChars = "~!@#%^*-_=+?"
	} else {
		specialChars = customChars[0]
	}
	return fmt.Sprintf("%s%s%s%d",
		acctest.RandStringFromCharSet(2, "ABCDEFGHIJKLMNOPQRSTUVWXZY"),
		acctest.RandStringFromCharSet(3, acctest.CharSetAlpha),
		acctest.RandStringFromCharSet(2, specialChars),
		acctest.RandIntRange(1000, 9999))
}

func ReplaceVarsForTest(rs *terraform.ResourceState, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{([[:word:]]+)}")

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return "replace_holder"
		}
		if rs != nil {
			if m == "id" {
				return rs.Primary.ID
			}
			v, ok := rs.Primary.Attributes[m]
			if ok {
				return v
			}
		}
		return ""
	}

	s := re.ReplaceAllStringFunc(linkTmpl, replaceFunc)
	return strings.Replace(s, "replace_holder/", "", 1), nil
}
