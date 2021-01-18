package provider

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/ssh"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"username": {
					Description: "A Git username. This is used to access a remote repository over HTTP. " +
						"It can also be sourced from the `GIT_USERNAME` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("GIT_USERNAME", nil),
				},
				"password": {
					Description: "A Git password. This is used to access a remote repository over HTTP. " +
						"It can also be sourced from the `GIT_PASSWORD` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("GIT_PASSWORD", nil),
				},
				"private_key": {
					Description: "A PEM-encoded private key. This is used to access a remote repository over SSH. " +
						"It can also be sourced from the `GIT_PRIVATE_KEY` environment variable.",
					Type:          schema.TypeString,
					Optional:      true,
					Sensitive:     true,
					DefaultFunc:   schema.EnvDefaultFunc("GIT_PRIVATE_KEY", nil),
					ConflictsWith: []string{"username", "password"},
				},
				"private_key_file": {
					Description: "A path to a PEM-encoded private key. This is used to access a remote repository over SSH. " +
						"It can also be sourced from the `GIT_PRIVATE_KEY_FILE` environment variable." +
						"Either this or `private_key` can be specified, but not both.",
					Type:          schema.TypeString,
					Optional:      true,
					DefaultFunc:   schema.EnvDefaultFunc("GIT_PRIVATE_KEY_FILE", nil),
					ConflictsWith: []string{"username", "password", "private_key"},
					Deprecated:    "Deprecated in favour of `private_key`",
				},
				"private_key_password": {
					Description: "An encryption password. Should be specified if the PEM-encoded private key contains a password " +
						"encrypted PEM block, otherwise password should be empty.",
					Type:          schema.TypeString,
					Optional:      true,
					DefaultFunc:   schema.EnvDefaultFunc("GIT_PRIVATE_KEY_PASSWORD", nil),
					ConflictsWith: []string{"username", "password"},
				},
				"ignore_host_key": {
					Description: "Set this to true to disable SSH host key verification. " +
						"This will accept any host key and is strongly discouraged. " +
						"It can also be sourced from the `GIT_IGNORE_HOST_KEY` environment variable.",
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("GIT_IGNORE_HOST_KEY", false),
				},
				"skip_tls_verify": {
					Description: "Set this to true to disable verification of the server's TLS certificate chain. " +
						"This is strongly discouraged. " +
						"It can also be sourced from the `GIT_SKIP_TLS_VERIFY` environment variable.",
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("GIT_SKIP_TLS_VERIFY", false),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"git_repository": dataSourceRepository(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type gitClient struct {
	auth transport.AuthMethod
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		gitClient := &gitClient{}

		username := d.Get("username").(string)
		password := d.Get("password").(string)
		privateKey := d.Get("private_key").(string)
		privateKeyFile := d.Get("private_key_file").(string)
		privateKeyPassword := d.Get("private_key_password").(string)

		if username != "" && password != "" {
			gitClient.auth = &githttp.BasicAuth{
				Username: username,
				Password: password,
			}
		}

		if privateKey != "" || privateKeyFile != "" {
			var publicKey *gitssh.PublicKeys
			var err error
			if privateKey != "" {
				publicKey, err = gitssh.NewPublicKeys(gitssh.DefaultUsername, []byte(privateKey), privateKeyPassword)
				if err != nil {
					return nil, diag.FromErr(err)
				}
			}
			if privateKeyFile != "" {
				publicKey, err = gitssh.NewPublicKeysFromFile(gitssh.DefaultUsername, privateKeyFile, privateKeyPassword)
				if err != nil {
					return nil, diag.FromErr(err)
				}
			}
			if v := d.Get("ignore_host_key").(bool); v {
				publicKey.HostKeyCallbackHelper = gitssh.HostKeyCallbackHelper{
					HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				}
			}

			gitClient.auth = publicKey
		}

		c := cleanhttp.DefaultClient()
		c.Transport = logging.NewTransport("Git", c.Transport)
		c.Transport = &userAgentTransport{
			userAgent: p.UserAgent("terraform-provider-git", version),
			transport: c.Transport,
		}
		if v, ok := d.GetOk("skip_tls_verify"); ok {
			c.Transport.(*http.Transport).TLSClientConfig = &tls.Config{
				InsecureSkipVerify: v.(bool),
			}
		}

		client.InstallProtocol("https", githttp.NewClient(c))

		return gitClient, nil
	}
}
