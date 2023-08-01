package acceptance

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ServiceFunc the resource query functions

// ResourceCheck resource check object
type ResourceCheck struct {
	resourceName    string
	resourceObject  interface{}
	getResourceFunc ServiceFunc
	resourceType    string
}

// CheckResourceDestroy check whether resources destroyed
func (rc *ResourceCheck) CheckResourceDestroy() resource.TestCheckFunc {
	if strings.Compare(rc.resourceType, dataSourceTypeCode) == 0 {
		return nil
	}

	return func(s *terraform.State) error {
		strs := strings.Split(rc.resourceName, ".")
		resourceType := strs[0]

		if resourceType == "" || resourceType == "data" {
			return fmt.Errorf("the format of the resource name is invalid, please check your configuration")
		}

		if rc.getResourceFunc == nil {
			return fmt.Errorf("the 'getResourceFunc' is nil, please set it during initialization")
		}

		conf := TestAccProvider.Meta().(*config.Config)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			if _, err := rc.getResourceFunc(conf, rs); err == nil {
				return fmt.Errorf("failed to destroy the %s resource: %s still exists",
					resourceType, rs.Primary.ID)
			}
		}
		return nil
	}
}

func (rc *ResourceCheck) checkResourceExists(s *terraform.State) error {
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
func (rc *ResourceCheck) CheckResourceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return rc.checkResourceExists(s)
	}
}

/*
CheckMultiResourcesExists checks whether multiple resources created by count are both existed.

	Parameters:
	  count: the expected number of resources that will be created.
*/
func (rc *ResourceCheck) CheckMultiResourcesExists(count int) resource.TestCheckFunc {
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
