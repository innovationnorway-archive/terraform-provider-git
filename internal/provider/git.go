package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	version "github.com/hashicorp/go-version"
)

type repoParams struct {
	URL   string
	Path  string
	Ref   plumbing.ReferenceName
	Auth  transport.AuthMethod
	Depth int
}

func getRepo(ctx context.Context, r repoParams) (*git.Repository, error) {
	if r.URL != "" {
		options := git.CloneOptions{
			URL:           r.URL,
			Auth:          r.Auth,
			ReferenceName: r.Ref,
			SingleBranch:  true,
			Depth:         r.Depth,
		}

		return git.CloneContext(ctx, memory.NewStorage(), memfs.New(), &options)
	}

	p, err := findGitDir(r.Path)
	if err != nil {
		return nil, err
	}

	return git.PlainOpen(p)
}

func getRef(repo *git.Repository, refName plumbing.ReferenceName) (*plumbing.Reference, error) {
	if refName != "" {
		return repo.Reference(refName, false)
	}

	return repo.Head()
}

func getRemoteURL(repo *git.Repository) string {
	remoteURL := ""
	remote, _ := repo.Remote(git.DefaultRemoteName)
	if remote != nil {
		remoteURL = remote.Config().URLs[0]
	}

	return remoteURL
}

func getTags(repo *git.Repository, ref *plumbing.Reference) ([]string, error) {
	var tags []string

	tagRefs, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	err = tagRefs.ForEach(func(tag *plumbing.Reference) error {
		if tag.Hash().String() == ref.Hash().String() {
			tags = append(tags, tag.Name().Short())
		}

		obj, _ := repo.TagObject(tag.Hash())
		if obj != nil {
			commit, _ := obj.Commit()
			if commit != nil {
				if commit.Hash.String() == ref.Hash().String() {
					tags = append(tags, tag.Name().Short())
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tags, nil
}

func getLatestTag(tags []string) string {
	var versions []*version.Version

	for _, t := range tags {
		v, _ := version.NewVersion(t)
		if v != nil {
			versions = append(versions, v)
		}
	}

	sort.Sort(sort.Reverse(version.Collection(versions)))
	if len(versions) > 0 {
		return versions[0].Original()
	}

	sort.Strings(tags)
	return tags[0]
}

func findGitDir(d string) (string, error) {
	if d == "" {
		d, _ = os.Getwd()
	}
	for d != string(os.PathSeparator) {
		p := filepath.Join(d, ".git")
		if stat, err := os.Stat(p); err != nil {
			if !os.IsNotExist(err) {
				return "", err
			}
		} else if stat.IsDir() {
			return d, nil
		}

		d = filepath.Dir(d)
	}
	return "", fmt.Errorf("not found: %s", d)
}
