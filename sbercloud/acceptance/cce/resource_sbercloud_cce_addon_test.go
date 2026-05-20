package cce

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/addons"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getAddonFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CceAddonV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE v3 client: %s", err)
	}
	return addons.Get(client, state.Primary.ID, state.Primary.Attributes["cluster_id"]).Extract()
}

func TestAccAddon_basic(t *testing.T) {
	var (
		addon addons.Addon

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "sbercloud_cce_addon.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&addon,
			getAddonFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddon_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id",
						"sbercloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "template_name", "metrics-server"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAddonImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAddonImportStateIdFunc(resName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterId, addonId string
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of CCE add-on is not found in the tfstate", resName)
		}
		clusterId = rs.Primary.Attributes["cluster_id"]
		addonId = rs.Primary.ID
		if clusterId == "" || addonId == "" {
			return "", fmt.Errorf("the CCE add-on ID is not exist or related CCE cluster ID is missing")
		}
		return fmt.Sprintf("%s/%s", clusterId, addonId), nil
	}
}

func testAccAddon_Base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cce_node" "test" {
  cluster_id        = sbercloud_cce_cluster.test.id
  name              = "%[2]s"
  flavor_id         = data.sbercloud_compute_flavors.test.ids[0]
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  key_pair          = sbercloud_kps_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
`, testAccNode_Base(rName), rName)
}

func testAccAddon_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cce_addon" "test" {
  cluster_id    = sbercloud_cce_cluster.test.id
  template_name = "metrics-server"
  depends_on    = [sbercloud_cce_node.test]
}
`, testAccAddon_Base(rName))
}

func TestAccAddon_values(t *testing.T) {
	var (
		addon addons.Addon

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "sbercloud_cce_addon.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&addon,
			getAddonFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddon_values_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id",
						"sbercloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.25.21"),
					resource.TestCheckResourceAttr(resourceName, "template_name", "autoscaler"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccAddon_values_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// the values not set, only check if the updating request is successful
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id",
						"sbercloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.25.21"),
					resource.TestCheckResourceAttr(resourceName, "template_name", "autoscaler"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

func testAccAddon_values_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cce_node_pool" "test" {
  cluster_id         = sbercloud_cce_cluster.test.id
  name               = "%[2]s"
  os                 = "CentOS 7.6"
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  initial_node_count = 4
  availability_zone  = data.sbercloud_availability_zones.test.names[0]
  key_pair           = sbercloud_kps_keypair.test.name
  scall_enable       = true
  min_node_count     = 2
  max_node_count     = 10
  priority           = 1
  type               = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}

data "sbercloud_cce_addon_template" "test" {
  cluster_id = sbercloud_cce_cluster.test.id
  name       = "autoscaler"
  version    = "1.25.21"
}
`, testAccNodePool_base(name), name)
}

func testAccAddon_values_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cce_addon" "test" {
  depends_on = [
    sbercloud_cce_node_pool.test,
  ]

  cluster_id    = sbercloud_cce_cluster.test.id
  template_name = "autoscaler"
  version       = "1.25.21"

  values {
    basic       = jsondecode(data.sbercloud_cce_addon_template.test.spec).basic
    custom_json = jsonencode(merge(
      jsondecode(data.sbercloud_cce_addon_template.test.spec).parameters.custom,
      {
        cluster_id = sbercloud_cce_cluster.test.id
        tenant_id  = "%[2]s"
        logLevel   = 3
      }
    ))
    flavor_json = jsonencode(jsondecode(data.sbercloud_cce_addon_template.test.spec).parameters.flavor1)
  }
}
`, testAccAddon_values_base(name), acceptance.SBC_PROJECT_ID)
}

func testAccAddon_values_step2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cce_addon" "test" {
  depends_on = [
    sbercloud_cce_node_pool.test,
  ]

  cluster_id    = sbercloud_cce_cluster.test.id
  template_name = "autoscaler"
  version       = "1.25.21"

  values {
    basic       = jsondecode(data.sbercloud_cce_addon_template.test.spec).basic
    custom_json = jsonencode(merge(
      jsondecode(data.sbercloud_cce_addon_template.test.spec).parameters.custom,
      {
        cluster_id = sbercloud_cce_cluster.test.id
        tenant_id  = "%[2]s"
        logLevel   = 4
      }
    ))
    flavor_json = jsonencode(jsondecode(data.sbercloud_cce_addon_template.test.spec).parameters.flavor2)
  }
}
`, testAccAddon_values_base(rName), acceptance.SBC_PROJECT_ID)
}

//import (
//	"fmt"
//	"testing"
//
//	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
//	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
//
//	"github.com/chnsz/golangsdk/openstack/cce/v3/addons"
//	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
//)
//
//func TestAccCCEAddonV3_basic(t *testing.T) {
//	var addon addons.Addon
//
//	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
//	resourceName := "sbercloud_cce_addon.test"
//	clusterName := "sbercloud_cce_cluster.test"
//
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
//		ProviderFactories: acceptance.TestAccProviderFactories,
//		CheckDestroy:      testAccCheckCCEAddonV3Destroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccCCEAddonV3_basic(rName),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckCCEAddonV3Exists(resourceName, clusterName, &addon),
//					resource.TestCheckResourceAttr(resourceName, "status", "running"),
//				),
//			},
//			{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//				ImportStateIdFunc: testAccCCEAddonImportStateIdFunc(),
//			},
//		},
//	})
//}
//
//func TestAccCCEAddonV3_values(t *testing.T) {
//	var addon addons.Addon
//
//	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
//	resourceName := "sbercloud_cce_addon.test"
//	clusterName := "sbercloud_cce_cluster.test"
//
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck: func() {
//			acceptance.TestAccPreCheck(t)
//		},
//		ProviderFactories: acceptance.TestAccProviderFactories,
//		CheckDestroy:      testAccCheckCCEAddonV3Destroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccCCEAddonV3_values(rName),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckCCEAddonV3Exists(resourceName, clusterName, &addon),
//					resource.TestCheckResourceAttr(resourceName, "status", "running"),
//				),
//			},
//		},
//	})
//}
//
//func testAccCheckCCEAddonV3Destroy(s *terraform.State) error {
//	config := acceptance.TestAccProvider.Meta().(*config.Config)
//	cceClient, err := config.CceAddonV3Client(acceptance.SBC_REGION_NAME)
//	if err != nil {
//		return fmtp.Errorf("Error creating SberCloud CCE Addon client: %s", err)
//	}
//
//	var clusterId string
//
//	for _, rs := range s.RootModule().Resources {
//		if rs.Type == "sbercloud_cce_cluster" {
//			clusterId = rs.Primary.ID
//		}
//
//		if rs.Type != "sbercloud_cce_addon" {
//			continue
//		}
//
//		if clusterId != "" {
//			_, err := addons.Get(cceClient, rs.Primary.ID, clusterId).Extract()
//			if err == nil {
//				return fmtp.Errorf("addon still exists")
//			}
//		}
//	}
//	return nil
//}
//
//func testAccCheckCCEAddonV3Exists(n string, cluster string, addon *addons.Addon) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return fmtp.Errorf("Not found: %s", n)
//		}
//		c, ok := s.RootModule().Resources[cluster]
//		if !ok {
//			return fmtp.Errorf("Cluster not found: %s", c)
//		}
//
//		if rs.Primary.ID == "" {
//			return fmtp.Errorf("No ID is set")
//		}
//		if c.Primary.ID == "" {
//			return fmtp.Errorf("Cluster id is not set")
//		}
//
//		config := acceptance.TestAccProvider.Meta().(*config.Config)
//		cceClient, err := config.CceAddonV3Client(acceptance.SBC_REGION_NAME)
//		if err != nil {
//			return fmtp.Errorf("Error creating SberCloud CCE Addon client: %s", err)
//		}
//
//		found, err := addons.Get(cceClient, rs.Primary.ID, c.Primary.ID).Extract()
//		if err != nil {
//			return err
//		}
//
//		if found.Metadata.Id != rs.Primary.ID {
//			return fmtp.Errorf("Addon not found")
//		}
//
//		*addon = *found
//
//		return nil
//	}
//}
//
//func testAccCCEAddonImportStateIdFunc() resource.ImportStateIdFunc {
//	return func(s *terraform.State) (string, error) {
//		var clusterID string
//		var addonID string
//		for _, rs := range s.RootModule().Resources {
//			if rs.Type == "sbercloud_cce_cluster" {
//				clusterID = rs.Primary.ID
//			} else if rs.Type == "sbercloud_cce_addon" {
//				addonID = rs.Primary.ID
//			}
//		}
//		if clusterID == "" || addonID == "" {
//			return "", fmtp.Errorf("resource not found: %s/%s", clusterID, addonID)
//		}
//		return fmt.Sprintf("%s/%s", clusterID, addonID), nil
//	}
//}
//
//func testAccCCEAddonV3_Base(rName string) string {
//	return fmt.Sprintf(`
//%s
//
//resource "sbercloud_cce_node" "test" {
//  cluster_id        = sbercloud_cce_cluster.test.id
//  name              = "%s"
//  flavor_id         = "c6nl.large.2"
//  availability_zone = data.sbercloud_availability_zones.test.names[0]
//  key_pair          = sbercloud_compute_keypair.test.name
//  os                = "CentOS 7.6"
//
//  root_volume {
//    size       = 50
//    volumetype = "SAS"
//  }
//  data_volumes {
//    size       = 100
//    volumetype = "SAS"
//  }
//}
//`, testAccCCENodeV3_Base(rName), rName)
//}
//
//func testAccCCEAddonV3_basic(rName string) string {
//	return fmt.Sprintf(`
//%s
//
//resource "sbercloud_cce_addon" "test" {
//  cluster_id    = sbercloud_cce_cluster.test.id
//  version       = "1.3.68"
//  template_name = "metrics-server"
//  depends_on    = [sbercloud_cce_node.test]
//}
//`, testAccCCEAddonV3_Base(rName))
//}
//
//func testAccCCEAddonV3_values(rName string) string {
//	return fmt.Sprintf(`
//%s
//
//resource "sbercloud_cce_node_pool" "test" {
//  cluster_id         = sbercloud_cce_cluster.test.id
//  name               = "%s"
//  os                 = "CentOS 7.6"
//  flavor_id          = "c6nl.large.2"
//  initial_node_count = 2
//  availability_zone  = data.sbercloud_availability_zones.test.names[0]
//  key_pair           = sbercloud_compute_keypair.test.name
//  scall_enable       = true
//  min_node_count     = 2
//  max_node_count     = 4
//  priority           = 1
//  type               = "vm"
//
//  root_volume {
//    size       = 50
//    volumetype = "SAS"
//  }
//  data_volumes {
//    size       = 100
//    volumetype = "SAS"
//  }
//}
//
//data "sbercloud_cce_addon_template" "test" {
//  cluster_id = sbercloud_cce_cluster.test.id
//  name       = "autoscaler"
//  version    = "1.30.18"
//}
//
//resource "sbercloud_cce_addon" "test" {
//  cluster_id    = sbercloud_cce_cluster.test.id
//  template_name = "autoscaler"
//  version       = "1.30.18"
//
//   values {
//    basic       = jsondecode(data.sbercloud_cce_addon_template.test.spec).basic
//    custom_json = jsonencode(merge(
//      jsondecode(data.sbercloud_cce_addon_template.test.spec).parameters.custom,
//      {
//        cluster_id = sbercloud_cce_cluster.test.id
//        tenant_id  = "%s"
//        logLevel   = 3
//      }
//    ))
//    flavor_json = jsonencode(jsondecode(data.sbercloud_cce_addon_template.test.spec).parameters.flavor1)
//  }
//
//  depends_on = [sbercloud_cce_node_pool.test]
//}
//`, testAccCCENodePool_Base(rName), rName, acceptance.SBC_PROJECT_ID)
//}
