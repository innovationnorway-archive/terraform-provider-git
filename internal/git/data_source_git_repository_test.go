package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func execGit(t *testing.T, arg ...string) string {
	t.Helper()

	output, err := exec.Command("git", arg...).Output()
	if err != nil {
		t.Fatal(err)
	}

	return strings.TrimSpace(string(output))
}

func TestDataSourceGitRepository_path(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(cwd, "..", ".git")

	branch := execGit(t, "rev-parse", "--abbrev-ref", "HEAD")
	commit := execGit(t, "rev-parse", "HEAD")

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGitRepositoryPathConfig(path),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "branch", strings.TrimSpace(branch)),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_sha", strings.TrimSpace(commit)),
				),
			},
		},
	})
}

func TestDataSourceGitRepository_branch(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(cwd, "..", ".git")

	branch := "master"
	commit := execGit(t, "rev-parse", branch)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGitRepositoryBranchConfig(path, branch),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "branch", branch),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_sha", strings.TrimSpace(commit)),
				),
			},
		},
	})
}

func TestDataSourceGitRepository_tag(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(cwd, "..", ".git")

	tag := "v0.1.0"
	commit := execGit(t, "rev-parse", tag)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGitRepositoryTagConfig(path, tag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "tag", tag),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_sha", strings.TrimSpace(commit)),
				),
			},
		},
	})
}

func TestDataSourceGitRepository_HTTPURL(t *testing.T) {
	url := "https://github.com/volcano-coffee-company/terraform-provider-git.git"
	tag := "v0.1.0"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGitRepositoryURLConfig(url, tag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "url", url),
					resource.TestCheckResourceAttr("data.git_repository.test", "tag", tag),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_sha", "396998df97b55acaa7d1645c0d90b3125ff51704"),
				),
			},
		},
	})
}

func testDataSourceGitRepositoryPathConfig(path string) string {
	return fmt.Sprintf(`
data git_repository "test" {
  path = "%s"
}
`, path)
}

func testDataSourceGitRepositoryBranchConfig(path string, branch string) string {
	return fmt.Sprintf(`
data git_repository "test" {
  path   = "%s"
  branch = "%s" 
}
`, path, branch)
}

func testDataSourceGitRepositoryTagConfig(path string, tag string) string {
	return fmt.Sprintf(`
data git_repository "test" {
  path = "%s"
  tag  = "%s" 
}
`, path, tag)
}

func testDataSourceGitRepositoryURLConfig(url string, tag string) string {
	return fmt.Sprintf(`
data git_repository "test" {
  url = "%s"
  tag  = "%s" 
}
`, url, tag)
}