package sbercloud

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	// Auth data
	SBC_REGION_NAME  = os.Getenv("SBC_REGION_NAME")
	SBC_ACCOUNT_NAME = os.Getenv("SBC_ACCOUNT_NAME")
	SBC_ACCESS_KEY   = os.Getenv("SBC_ACCESS_KEY")
	SBC_SECRET_KEY   = os.Getenv("SBC_SECRET_KEY")
	// Test data
	TEST_SBC_AVAILABILITY_ZONE = os.Getenv("TEST_SBC_AVAILABILITY_ZONE")
	TEST_SBC_PROJECT_ID        = os.Getenv("TEST_SBC_PROJECT_ID")
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"sbercloud": testAccProvider,
	}
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	if SBC_REGION_NAME == "" {
		t.Fatal("SBC_REGION_NAME must be set for acceptance tests")
	}
	if SBC_ACCOUNT_NAME == "" {
		t.Fatal("SBC_ACCOUNT_NAME must be set for acceptance tests")
	}
	if SBC_ACCESS_KEY == "" {
		t.Fatal("SBC_ACCESS_KEY must be set for acceptance tests")
	}
	if SBC_SECRET_KEY == "" {
		t.Fatal("SBC_SECRET_KEY must be set for acceptance tests")
	}
	if TEST_SBC_AVAILABILITY_ZONE == "" {
		t.Fatal("TEST_SBC_AVAILABILITY_ZONE must be set for acceptance tests")
	}
	if TEST_SBC_PROJECT_ID == "" {
		t.Fatal("TEST_SBC_PROJECT_ID must be set for acceptance tests")
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

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
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
