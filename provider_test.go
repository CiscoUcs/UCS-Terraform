package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"ucs": testAccProvider,
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

func TestproviderConfigure(t *testing.T) {
	resourceData := &schema.ResourceData{}
	client, err := providerConfigure(resourceData)
	if err != nil {
		t.Errorf("providerConfigure() returned unexpected error %v", err)
	}

	if client == nil {
		t.Errorf("providerConfigure() returned nil client; expected an instance of a UCSClient")
	}
}
