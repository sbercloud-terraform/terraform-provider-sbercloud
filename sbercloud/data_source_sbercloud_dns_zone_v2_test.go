package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestDNSZoneV2DataSourcePublic(t *testing.T) {
	zoneName := fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))
	dataSourceName := "data.sbercloud_dns_zone.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDNSZoneV2SourceDNSBasePublic(zoneName),
			},
			{
				Config: testDNSZoneV2SourceDNSByNamePublic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", zoneName),
					resource.TestCheckResourceAttr(dataSourceName, "zone_type", "public"),
					testAccCheckDmsZoneV2DataSourceID(dataSourceName),
				),
			},
		},
	})
}

func TestDNSZoneV2DataSourcePrivate(t *testing.T) {
	zoneName := fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))
	vpcName := fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
	dataSourceName := "data.sbercloud_dns_zone.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDNSZoneV2SourceDNSByNamePrivate(zoneName, vpcName),
			},
			{
				Config: testDNSZoneV2SourceDNSByNamePrivate(zoneName, vpcName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", zoneName),
					resource.TestCheckResourceAttr(dataSourceName, "zone_type", "private"),
					testAccCheckDmsZoneV2DataSourceID(dataSourceName),
				),
			},
		},
	})
}

func testDNSZoneV2SourceDNSBasePublic(zoneName string) string {
	return fmt.Sprintf(`
resource "sbercloud_dns_zone" "test" {
  name = "%s"
  zone_type = "public"
}
`, zoneName)
}

func testDNSZoneV2SourceDNSByNamePublic(zoneName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_dns_zone" "test" {
  name = sbercloud_dns_zone.test.name
  zone_type = "public"
}
`, testDNSZoneV2SourceDNSBasePublic(zoneName))
}

func testDNSZoneV2SourceDNSBasePrivate(zoneName, vpcName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name                  = "%s"
  cidr                  = "192.168.0.0/16"
}
resource "sbercloud_dns_zone" "test" {
  name = "%s"
  zone_type = "private"
  router {
	router_id = sbercloud_vpc.test.id
  }
}
`, vpcName, zoneName)
}

func testDNSZoneV2SourceDNSByNamePrivate(zoneName, vpcName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_dns_zone" "test" {
  name = sbercloud_dns_zone.test.name
  zone_type = "private"
}
`, testDNSZoneV2SourceDNSBasePrivate(zoneName, vpcName))
}

func testAccCheckDmsZoneV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DNS zone data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DNS zone data source ID not set")
		}

		return nil
	}
}
