package sbercloud

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/recordsets"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func randomZoneName() string {
	// TODO: why does back-end convert name to lowercase?
	return fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))
}

func parseDNSV2RecordSetID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 {
		return "", "", fmt.Errorf("Unable to determine DNS record set ID from raw ID: %s", id)
	}

	zoneID := idParts[0]
	recordsetID := idParts[1]

	return zoneID, recordsetID, nil
}

func TestAccDNSV2RecordSet_basic(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("sbercloud_dns_recordset.recordset_1", &recordset),
					resource.TestCheckResourceAttr(
						"sbercloud_dns_recordset.recordset_1", "description", "a record set"),
					resource.TestCheckResourceAttr(
						"sbercloud_dns_recordset.recordset_1", "records.0", "10.1.0.0"),
				),
			},
			{
				ResourceName:      "sbercloud_dns_recordset.recordset_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDNSV2RecordSet_update(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sbercloud_dns_recordset.recordset_1", "name", zoneName),
					resource.TestCheckResourceAttr("sbercloud_dns_recordset.recordset_1", "ttl", "6000"),
					resource.TestCheckResourceAttr("sbercloud_dns_recordset.recordset_1", "type", "A"),
					resource.TestCheckResourceAttr(
						"sbercloud_dns_recordset.recordset_1", "description", "an updated record set"),
					resource.TestCheckResourceAttr(
						"sbercloud_dns_recordset.recordset_1", "records.0", "10.1.0.1"),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_readTTL(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_readTTL(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("sbercloud_dns_recordset.recordset_1", &recordset),
					resource.TestMatchResourceAttr(
						"sbercloud_dns_recordset.recordset_1", "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func testAccCheckDNSV2RecordSetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dnsClient, err := config.DnsV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_dns_recordset" {
			continue
		}

		zoneID, recordsetID, err := parseDNSV2RecordSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err == nil {
			return fmt.Errorf("Record set still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2RecordSetExists(n string, recordset *recordsets.RecordSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		dnsClient, err := config.DnsV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
		}

		zoneID, recordsetID, err := parseDNSV2RecordSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err != nil {
			return err
		}

		if found.ID != recordsetID {
			return fmt.Errorf("Record set not found")
		}

		*recordset = *found

		return nil
	}
}

func testAccDNSV2RecordSet_basic(zoneName string) string {
	return fmt.Sprintf(`
		resource "sbercloud_dns_zone" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "a zone"
			ttl = 6000
			#type = "PRIMARY"
		}

		resource "sbercloud_dns_recordset" "recordset_1" {
			zone_id = sbercloud_dns_zone.zone_1.id
			name = "%s"
			type = "A"
			description = "a record set"
			ttl = 3000
			records = ["10.1.0.0"]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSet_update(zoneName string) string {
	return fmt.Sprintf(`
		resource "sbercloud_dns_zone" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "an updated zone"
			ttl = 6000
			#type = "PRIMARY"
		}

		resource "sbercloud_dns_recordset" "recordset_1" {
			zone_id = sbercloud_dns_zone.zone_1.id
			name = "%s"
			type = "A"
			description = "an updated record set"
			ttl = 6000
			records = ["10.1.0.1"]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSet_readTTL(zoneName string) string {
	return fmt.Sprintf(`
		resource "sbercloud_dns_zone" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "an updated zone"
			ttl = 6000
			#type = "PRIMARY"
		}

		resource "sbercloud_dns_recordset" "recordset_1" {
			zone_id = sbercloud_dns_zone.zone_1.id
			name = "%s"
			type = "A"
			records = ["10.1.0.2"]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSet_timeout(zoneName string) string {
	return fmt.Sprintf(`
		resource "sbercloud_dns_zone" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "an updated zone"
			ttl = 6000
			#type = "PRIMARY"
		}

		resource "sbercloud_dns_recordset" "recordset_1" {
			zone_id = sbercloud_dns_zone.zone_1.id
			name = "%s"
			type = "A"
			ttl = 3000
			records = ["10.1.0.3"]

			timeouts {
				create = "5m"
				update = "5m"
				delete = "5m"
			}
		}
	`, zoneName, zoneName)
}
