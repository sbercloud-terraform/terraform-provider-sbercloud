package sbercloud

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/pathorcontents"
)

var (
	SBC_ACCESS_KEY                 = os.Getenv("SBC_ACCESS_KEY")
	SBC_ACCOUNT_NAME               = os.Getenv("SBC_ACCOUNT_NAME")
	SBC_ADMIN                      = os.Getenv("SBC_ADMIN")
	SBC_DOMAIN_ID                  = os.Getenv("SBC_DOMAIN_ID")
	SBC_DOMAIN_NAME                = os.Getenv("SBC_DOMAIN_NAME")
	SBC_ENTERPRISE_PROJECT_ID_TEST = os.Getenv("SBC_ENTERPRISE_PROJECT_ID_TEST")
	SBC_PROJECT_ID                 = os.Getenv("SBC_PROJECT_ID")
	SBC_REGION_NAME                = os.Getenv("SBC_REGION_NAME")
	SBC_SECRET_KEY                 = os.Getenv("SBC_SECRET_KEY")
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"sbercloud": testAccProvider,
	}
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	if SBC_REGION_NAME == "" {
		t.Fatal("SBC_REGION_NAME must be set for acceptance tests")
	}
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
}

func testAccPreCheckAdminOnly(t *testing.T) {
	if SBC_ADMIN == "" {
		t.Skip("SBC_ADMIN must be set for acceptance tests")
	}
}

func testAccPreCheckEpsID(t *testing.T) {
	if SBC_ENTERPRISE_PROJECT_ID_TEST == "" {
		t.Skip("This environment does not support EPS_ID tests")
	}
}

func testAccPreCheckOBS(t *testing.T) {
	if SBC_ACCESS_KEY == "" || SBC_SECRET_KEY == "" {
		t.Skip("SBC_ACCESS_KEY and SBC_SECRET_KEY must be set for OBS acceptance tests")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func envVarContents(varName string) (string, error) {
	contents, _, err := pathorcontents.Read(os.Getenv(varName))
	if err != nil {
		return "", fmt.Errorf("Error reading %s: %s", varName, err)
	}
	return contents, nil
}

func envVarFile(varName string) (string, error) {
	contents, err := envVarContents(varName)
	if err != nil {
		return "", err
	}

	tmpFile, err := ioutil.TempFile("", varName)
	if err != nil {
		return "", fmt.Errorf("Error creating temp file: %s", err)
	}
	if _, err := tmpFile.Write([]byte(contents)); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error writing temp file: %s", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error closing temp file: %s", err)
	}
	return tmpFile.Name(), nil
}
