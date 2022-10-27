package acceptance

import (
	"encoding/json"
	"fmt"
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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud"
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
func (rc *resourceCheck) CheckResourceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rc.resourceName]
		if !ok {
			return fmtp.Errorf("Can not found the resource or data source in state: %s", rc.resourceName)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("No id set for the resource or data source: %s", rc.resourceName)
		}
		if strings.EqualFold(rc.resourceType, dataSourceTypeCode) {
			return nil
		}

		if rc.getResourceFunc != nil {
			conf := TestAccProvider.Meta().(*config.Config)
			r, err := rc.getResourceFunc(conf, rs)
			if err != nil {
				return fmtp.Errorf("checking resource %s %s exists error: %s ",
					rc.resourceName, rs.Primary.ID, err)
			}
			if rc.resourceObject != nil {
				b, err := json.Marshal(r)
				if err != nil {
					return fmtp.Errorf("marshaling resource %s %s error: %s ",
						rc.resourceName, rs.Primary.ID, err)
				}
				json.Unmarshal(b, rc.resourceObject)
			} else {
				logp.Printf("[WARN] The 'resourceObject' is nil, please set it during initialization.")
			}
		} else {
			return fmtp.Errorf("The 'getResourceFunc' is nil, please set it.")
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
