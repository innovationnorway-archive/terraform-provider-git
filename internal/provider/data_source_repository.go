package provider

import (
	"context"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceRepository() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to retrieve information about a Git repository.",

		ReadContext: dataSourceRepositoryRead,

		Schema: map[string]*schema.Schema{
			"path": {
				Description:      "The repository path. It can also be sourced from the `GIT_DIR` environment variable.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				DefaultFunc:      schema.EnvDefaultFunc("GIT_DIR", nil),
			},

			"url": {
				Description:      "The repository URL.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				ExactlyOneOf:     []string{"path", "url"},
			},

			"branch": {
				Description:      "The name of the branch.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},

			"tag": {
				Description:      "The name of the tag.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				ConflictsWith:    []string{"branch"},
			},

			"commit_sha": {
				Description: "The SHA-1 hash of the current commit.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	auth := meta.(*gitClient).auth

	params := repoParams{
		Auth:  auth,
		Depth: 1,
	}

	if v, ok := d.GetOk("path"); ok {
		params.Path = v.(string)
	}

	if v, ok := d.GetOk("url"); ok {
		params.URL = v.(string)
	}

	if v, ok := d.GetOk("branch"); ok {
		params.Ref = plumbing.NewBranchReferenceName(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		params.Ref = plumbing.NewTagReferenceName(v.(string))
	}

	repo, err := getRepo(ctx, params)
	if err != nil {
		return diag.Errorf("error getting repository: %s", err)
	}

	ref, err := getRef(repo, params.Ref)
	if err != nil {
		return diag.Errorf("error getting reference: %s", err)
	}

	if params.Path == "" {
		d.Set("path", params.Path)
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
			return diag.Errorf("errr getting tags: %s", err)
		}

		if tags != nil {
			d.Set("tag", getLatestTag(tags))
		}
	}

	d.Set("commit_sha", ref.Hash().String())
	d.SetId(ref.Hash().String())

	return nil
}
