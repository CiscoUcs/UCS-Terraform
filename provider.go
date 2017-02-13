package main

import (
	"github.com/CiscoCloud/terraform-provider-ucs/ucsclient"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ip_address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "UCS Manager IP address or CIMC IP address.",
			},

			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user's name to access the UCS Management.",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The password to access the UCS Management.",
			},

			"tslinsecureskipverify": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The TSL insecure skip verify",
			},

			"log_level": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The log level",
			},

			"log_filename": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The log filename",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"ucs_service_profile": resourceUcsServiceProfile(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := ucsclient.Config{
		AppName:               "UCS",
		IpAddress:             d.Get("ip_address").(string),
		Username:              d.Get("username").(string),
		Password:              d.Get("password").(string),
		TslInsecureSkipVerify: d.Get("tslinsecureskipverify").(bool),
		LogFilename:           d.Get("log_filename").(string),
		LogLevel:              d.Get("log_level").(int),
	}

	return config.Client(), nil
}
