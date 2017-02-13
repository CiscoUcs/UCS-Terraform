package main

import (
	"testing"
)

func TestValidateCIDR(t *testing.T) {
	cidr := "1.2.3.4/8"
	err := validateCIDR(cidr)
	if err != nil {
		t.Errorf("CIDR %s returned error: %s", cidr, err)
	}

	cidr = "hola"
	err = validateCIDR(cidr)
	if err == nil {
		t.Errorf(`Error expected but got nil with cidr = "%s"`, cidr)
	}
}
