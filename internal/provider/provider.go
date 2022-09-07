package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GIT_USERNAME", nil),
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("GIT_PASSWORD", nil),
			},

			"private_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("GIT_PRIVATE_KEY", nil),
				ConflictsWith: []string{"username", "password"},
			},

			"private_key_file": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("GIT_PRIVATE_KEY_FILE", nil),
				ConflictsWith: []string{"username", "password", "private_key"},
				Deprecated:    "Deprecated in favour of `private_key`",
			},

			"ignore_host_key": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GIT_IGNORE_HOST_KEY", false),
			},

			"skip_tls_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GIT_SKIP_TLS_VERIFY", false),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"git_repository": dataSourceGitRepository(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		Username:              d.Get("username").(string),
		Password:              d.Get("password").(string),
		PrivateKey:            d.Get("private_key").(string),
		InsecureIgnoreHostKey: d.Get("ignore_host_key").(bool),
		InsecureSkipTLSVerify: d.Get("skip_tls_verify").(bool),
	}

	return config.Client()
}
