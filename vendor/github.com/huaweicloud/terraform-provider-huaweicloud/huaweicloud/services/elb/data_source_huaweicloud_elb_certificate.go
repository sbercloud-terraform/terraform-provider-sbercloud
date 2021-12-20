/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package elb

import (
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/elb/v3/certificates"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// DataSourceELBCertificateV3 the data source of "huaweicloud_elb_certificate"
// Dedicated ELB
func DataSourceELBCertificateV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceELbCertificateV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceELbCertificateV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud Dedicated ELB(V3) Client: %s", err)
	}

	listOpts := certificates.ListOpts{
		Name: d.Get("name").(string),
	}
	r, err := certificates.List(client, listOpts).AllPages()
	certs, err := certificates.ExtractCertificates(r)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve certs from Dedicated ELB(V3): %s", err)
	}
	logp.Printf("[DEBUG] Get certificate list: %#v", certs)

	if len(certs) > 0 {
		setCertificateAttributes(d, certs[0])
		if err != nil {
			return err
		}
	} else {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	return nil
}

func setCertificateAttributes(d *schema.ResourceData, c certificates.Certificate) error {
	d.SetId(c.ID)

	var expiration string
	tm, err := time.Parse(time.RFC3339, c.ExpireTime)
	if err != nil {
		// If the format of ExpireTime is not expected, set the original value directly.
		expiration = c.ExpireTime
		logp.Printf("[WAIN] The format of the ExpireTime field of the Dedicated ELB certificate "+
			"is not expected:%s", c.ExpireTime)
	} else {
		expiration = tm.Format("2006-01-02 15:04:05 MST")
	}

	mErr := multierror.Append(nil,
		d.Set("name", c.Name),
		d.Set("domain", c.Domain),
		d.Set("description", c.Description),
		d.Set("type", c.Type),
		d.Set("expiration", expiration),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("error setting Dedicated ELB Certificate fields: %s", err)
	}
	return nil
}
