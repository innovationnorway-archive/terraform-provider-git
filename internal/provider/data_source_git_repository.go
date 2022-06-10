package provider

import (
	"context"

	"github.com/go-git/go-git/v5/plumbing"
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
				ExactlyOneOf: []string{"path", "url"},
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

	repo, err := getRepo(params)
	if err != nil {
		return diag.Errorf("unable to get repository: %s", err)
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

	d.Set("commit_sha", ref.Hash().String())
	d.SetId(ref.Name().String())

	return nil
}
