package provider

import (
	"context"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGitRepository() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGitRepositoryRead,

		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				DefaultFunc:  schema.EnvDefaultFunc("GIT_DIR", nil),
			},

			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"branch": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tag": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"branch"},
			},

			"commit_sha": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"relative_path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"clean": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceGitRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)

	params := RepoParams{
		Auth: m.Auth,
	}

	if v, ok := d.GetOk("url"); ok {
		params.URL = v.(string)
	}

	if v, ok := d.GetOk("path"); ok {
		params.Path = v.(string)
	}

	if v, ok := d.GetOk("branch"); ok {
		params.Ref = plumbing.NewBranchReferenceName(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		params.Ref = plumbing.NewTagReferenceName(v.(string))
	}

	// Try to find the repository in current direcrory and above if not found
	var repo *git.Repository
	if params.Path == "" && params.URL == "" {
		path, err := os.Getwd()
		initialPath := path
		if err != nil {
			return diag.Errorf("failed to get current directory: %s", err)
		}

		for {
			if path == "/" {
				return diag.Errorf("unable to find repository: %s", err)
			}
			params.Path = path

			repo, err = getRepo(params)
			if err != nil {
				path = filepath.Dir(path)
				continue
			} else {
				v, err := filepath.Rel(path, initialPath)
				if err != nil {
					tflog.Debug(ctx, "error relative_path", map[string]interface{}{"relative_path": v})
				} else {
					d.Set("relative_path", v)
				}
				break
			}
		}
	} else {
		var err error
		repo, err = getRepo(params)
		if err != nil {
			return diag.Errorf("unable to find repository: %s", err)
		}
	}

	ref, err := getRef(repo, params.Ref)
	if err != nil {
		return diag.Errorf("unable to get reference: %s", err)
	}

	if params.URL == "" {
		d.Set("url", getRemoteURL(repo))
	}

	if ref.Name().IsBranch() {
		d.Set("branch", ref.Name().Short())
	}

	if ref.Name().IsTag() {
		d.Set("tag", ref.Name().Short())
	} else {
		tags, err := getTags(repo, ref)
		if err != nil {
			return diag.Errorf("unable to get tags: %s", err)
		}

		if tags != nil {
			d.Set("tag", getLatestTag(tags))
		}
	}

	clean, err := IsClean(repo)
	if err != nil {
		return diag.Errorf("unable to get status: %s", err)
	}
	d.Set("clean", clean)

	d.Set("commit_sha", ref.Hash().String())
	d.SetId(ref.Name().String())

	return nil
}
