package sbercloud

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccSberCloudObsBucketObjectDataSource_content(t *testing.T) {
	rInt := acctest.RandInt()
	resourceConf, dataSourceConf := testAccSberCloudObsBucketObjectDataSource_content(rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists("sbercloud_obs_bucket_object.object"),
				),
			},
			{
				Config: dataSourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsObsObjectDataSourceExists("data.sbercloud_obs_bucket_object.obj"),
					resource.TestCheckResourceAttr("data.sbercloud_obs_bucket_object.obj", "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr("data.sbercloud_obs_bucket_object.obj", "storage_class", "STANDARD"),
				),
			},
		},
	})
}

func TestAccSberCloudObsBucketObjectDataSource_source(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-acc-obs-obj-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	rInt := acctest.RandInt()

	// write test data to the tempfile
	for i := 0; i < 1024; i++ {
		_, err := tmpFile.WriteString("test obs object file storage")
		if err != nil {
			t.Fatal(err)
		}
	}
	tmpFile.Close()

	resourceConf, dataSourceConf := testAccSberCloudObsBucketObjectDataSource_source(rInt, tmpFile.Name())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists("sbercloud_obs_bucket_object.object"),
				),
			},
			{
				Config: dataSourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsObsObjectDataSourceExists("data.sbercloud_obs_bucket_object.obj"),
					resource.TestCheckResourceAttr("data.sbercloud_obs_bucket_object.obj", "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr("data.sbercloud_obs_bucket_object.obj", "storage_class", "STANDARD"),
				),
			},
		},
	})
}

func TestAccSberCloudObsBucketObjectDataSource_allParams(t *testing.T) {
	rInt := acctest.RandInt()
	resourceConf, dataSourceConf := testAccSberCloudObsBucketObjectDataSource_allParams(rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists("sbercloud_obs_bucket_object.object"),
				),
			},
			{
				Config: dataSourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsObsObjectDataSourceExists("data.sbercloud_obs_bucket_object.obj"),
					resource.TestCheckResourceAttr("data.sbercloud_obs_bucket_object.obj", "content_type", "application/unknown"),
					resource.TestCheckResourceAttr("data.sbercloud_obs_bucket_object.obj", "storage_class", "STANDARD"),
				),
			},
		},
	})
}

func testAccCheckAwsObsObjectDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Obs object data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Obs object data source ID not set")
		}

		bucket := rs.Primary.Attributes["bucket"]
		key := rs.Primary.Attributes["key"]

		config := testAccProvider.Meta().(*config.Config)
		obsClient, err := config.ObjectStorageClient(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud OBS client: %s", err)
		}

		respList, err := obsClient.ListObjects(&obs.ListObjectsInput{
			Bucket: bucket,
			ListObjsInput: obs.ListObjsInput{
				Prefix: key,
			},
		})
		if err != nil {
			return getObsError("Error listing objects of OBS bucket", bucket, err)
		}

		var exist bool
		for _, content := range respList.Contents {
			if key == content.Key {
				exist = true
				break
			}
		}
		if !exist {
			return fmt.Errorf("object %s not found in bucket %s", key, bucket)
		}

		return nil
	}
}

func testAccSberCloudObsBucketObjectDataSource_content(randInt int) (string, string) {
	resource := fmt.Sprintf(`
resource "sbercloud_obs_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
}
resource "sbercloud_obs_bucket_object" "object" {
	bucket = sbercloud_obs_bucket.object_bucket.bucket
	key = "test-key-%d"
	content = "some_bucket_content"
}
`, randInt, randInt)

	dataSource := fmt.Sprintf(`%s
data "sbercloud_obs_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%d"
	key = "test-key-%d"
}`, resource, randInt, randInt)

	return resource, dataSource
}

func testAccSberCloudObsBucketObjectDataSource_source(randInt int, source string) (string, string) {
	resource := fmt.Sprintf(`
resource "sbercloud_obs_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
}
resource "sbercloud_obs_bucket_object" "object" {
	bucket = sbercloud_obs_bucket.object_bucket.bucket
	key = "test-key-%d"
	source = "%s"
	content_type = "binary/octet-stream"
}
`, randInt, randInt, source)

	dataSource := fmt.Sprintf(`%s
data "sbercloud_obs_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%d"
	key = "test-key-%d"
}`, resource, randInt, randInt)

	return resource, dataSource
}

func testAccSberCloudObsBucketObjectDataSource_allParams(randInt int) (string, string) {
	resource := fmt.Sprintf(`
resource "sbercloud_obs_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
}
resource "sbercloud_obs_bucket_object" "object" {
	bucket = sbercloud_obs_bucket.object_bucket.bucket
	key = "test-key-%d"
	content = <<CONTENT
	{"msg": "Hi there!"}
CONTENT
	acl = "private"
	content_type = "application/unknown"
	storage_class = "STANDARD"
}
`, randInt, randInt)

	dataSource := fmt.Sprintf(`%s
data "sbercloud_obs_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%d"
	key = "test-key-%d"
}`, resource, randInt, randInt)

	return resource, dataSource
}
