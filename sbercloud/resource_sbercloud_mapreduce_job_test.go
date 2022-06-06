package sbercloud

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/mrs/v2/jobs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	mrsRes "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mrs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccMrsMapReduceJob_basic(t *testing.T) {
	var job jobs.Job
	resourceName := "sbercloud_mapreduce_job.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	pwd := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMRSV2JobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceJobConfig_basic(rName, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2JobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", mrsRes.JobHiveSQL),
					resource.TestCheckResourceAttr(resourceName, "sql", "SHOW DATABASES;"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMRSClusterSubResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCheckMRSV2JobDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.MrsV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sbercloud mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_mapreduce_job" {
			continue
		}

		_, err := jobs.Get(client, rs.Primary.Attributes["cluster_id"], rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("MRS cluster (%s) is still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMRSV2JobExists(n string, job *jobs.Job) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No MRS cluster ID")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.MrsV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sbercloud MRS client: %s ", err)
		}

		found, err := jobs.Get(client, rs.Primary.Attributes["cluster_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*job = *found
		return nil
	}
}

func testAccMRSClusterSubResourceImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["cluster_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID), nil
	}
}

func testAccMrsMapReduceJobConfig_base(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_mapreduce_cluster" "test" {
  availability_zone  = data.sbercloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "ANALYSIS"
  version            = "MRS 2.1.0"
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = sbercloud_vpc_subnet.test.id
  vpc_id             = sbercloud_vpc.test.id
  component_list     = ["Hadoop", "Spark", "Hive", "Tez"]
  safe_mode          = false

  master_nodes {
    flavor            = "c6.xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 100
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "c6.xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 100
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd)
}

func testAccMrsMapReduceJobConfig_basic(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_mapreduce_job" "test" {
  cluster_id   = sbercloud_mapreduce_cluster.test.id
  name         = "%s"
  type         = "HiveSql"
  sql          = "SHOW DATABASES;"
}`, testAccMrsMapReduceJobConfig_base(rName, pwd), rName)
}
