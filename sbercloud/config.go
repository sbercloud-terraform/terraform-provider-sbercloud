package sbercloud

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/huaweicloud/golangsdk"
	huaweisdk "github.com/huaweicloud/golangsdk/openstack"
)

const (
	serviceProjectLevel string = "project"
	serviceDomainLevel  string = "domain"
)

type Config struct {
	AccessKey        string
	SecretKey        string
	AccountName      string
	IdentityEndpoint string
	Insecure         bool
	Password         string
	Region           string
	ProjectName      string
	Username         string
	terraformVersion string

	HwClient *golangsdk.ProviderClient

	DomainClient *golangsdk.ProviderClient
}

func (c *Config) LoadAndValidate() error {
	err := fmt.Errorf("Must config aksk or username password to be authorized")
	if c.Password != "" {
		err = buildClientByPassword(c)
	} else if c.AccessKey != "" {
		err = buildClientByAKSK(c)

	}
	if err != nil {
		return err
	}

	return nil
}

func generateTLSConfig(c *Config) (*tls.Config, error) {
	config := &tls.Config{}
	if c.Insecure {
		config.InsecureSkipVerify = true
	}
	return config, nil
}

func genClient(c *Config, ao golangsdk.AuthOptionsProvider) (*golangsdk.ProviderClient, error) {
	client, err := huaweisdk.NewClient(ao.GetIdentityEndpoint())
	if err != nil {
		return nil, err
	}

	// Set UserAgent
	client.UserAgent.Prepend(httpclient.TerraformUserAgent(c.terraformVersion))

	config, err := generateTLSConfig(c)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}

	client.HTTPClient = http.Client{
		Transport: &LogRoundTripper{
			Rt:      transport,
			OsDebug: logging.IsDebugOrHigher(),
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if client.AKSKAuthOptions.AccessKey != "" {
				golangsdk.ReSign(req, golangsdk.SignOptions{
					AccessKey: client.AKSKAuthOptions.AccessKey,
					SecretKey: client.AKSKAuthOptions.SecretKey,
				})
			}
			return nil
		},
	}

	// Validate authentication normally.
	err = huaweisdk.Authenticate(client, ao)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func buildClientByAKSK(c *Config) error {
	var pao, dao golangsdk.AKSKAuthOptions

	pao = golangsdk.AKSKAuthOptions{
		ProjectName: c.ProjectName,
	}

	dao = golangsdk.AKSKAuthOptions{
		Domain: c.AccountName,
	}

	for _, ao := range []*golangsdk.AKSKAuthOptions{&pao, &dao} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.AccessKey = c.AccessKey
		ao.SecretKey = c.SecretKey
	}
	return genClients(c, pao, dao)
}

func buildClientByPassword(c *Config) error {
	var pao, dao golangsdk.AuthOptions

	pao = golangsdk.AuthOptions{
		DomainName: c.AccountName,
		TenantName: c.ProjectName,
	}

	dao = golangsdk.AuthOptions{
		DomainName: c.AccountName,
	}

	for _, ao := range []*golangsdk.AuthOptions{&pao, &dao} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.Password = c.Password
		ao.Username = c.Username
	}
	return genClients(c, pao, dao)
}

func genClients(c *Config, pao, dao golangsdk.AuthOptionsProvider) error {
	client, err := genClient(c, pao)
	if err != nil {
		return err
	}
	c.HwClient = client

	client, err = genClient(c, dao)
	if err == nil {
		c.DomainClient = client
	}
	return err
}

func (c *Config) determineRegion(region string) string {
	// If a resource-level region was not specified, and a provider-level region was set,
	// use the provider-level region.
	if region == "" && c.Region != "" {
		region = c.Region
	}

	log.Printf("[DEBUG] SberCloud Region is: %s", region)
	return region
}

func (c *Config) getHwEndpointType() golangsdk.Availability {
	return golangsdk.AvailabilityPublic
}

func (c *Config) identityV3Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewIdentityV3(c.DomainClient, golangsdk.EndpointOpts{
		//Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}
