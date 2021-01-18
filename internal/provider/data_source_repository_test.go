package provider

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRepository_URL(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryURLConfig("https://github.com/innovationnorway/terraform-module-acctest.git"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "branch", "master"),
					resource.TestCheckResourceAttr("data.git_repository.test", "tag", "v0.1.0"),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_sha", "a48e01ced30c807293088746dc2b6715d78882f1"),
				),
			},
		},
	})
}

func TestAccDataSourceRepository_Path(t *testing.T) {
	dir, err := ioutil.TempDir("", "acctest-*")
	if err != nil {
		t.Fatal(err)
	}
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: "https://github.com/innovationnorway/terraform-module-acctest.git",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryPathConfig(filepath.ToSlash(dir)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "url", "https://github.com/innovationnorway/terraform-module-acctest.git"),
					resource.TestCheckResourceAttr("data.git_repository.test", "branch", "master"),
					resource.TestCheckResourceAttr("data.git_repository.test", "tag", "v0.1.0"),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_sha", "a48e01ced30c807293088746dc2b6715d78882f1"),
				),
			},
		},
	})
}

func TestAccDataSourceRepository_SubPath(t *testing.T) {
	dir, err := ioutil.TempDir("", "acctest-*")
	if err != nil {
		t.Fatal(err)
	}
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: "https://github.com/innovationnorway/terraform-module-acctest.git",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryPathConfig(filepath.ToSlash(filepath.Join(dir, "test"))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "url", "https://github.com/innovationnorway/terraform-module-acctest.git"),
					resource.TestCheckResourceAttr("data.git_repository.test", "branch", "master"),
					resource.TestCheckResourceAttr("data.git_repository.test", "tag", "v0.1.0"),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_sha", "a48e01ced30c807293088746dc2b6715d78882f1"),
				),
			},
		},
	})
}

func TestAccDataSourceRepository_Tag(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryTagConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "url", "https://github.com/innovationnorway/terraform-module-acctest.git"),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_sha", "ea59fdb83b2ad18aa1da930f99aa34a1f325c009"),
				),
			},
		},
	})
}

func TestAccDataSourceRepository_Branch(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryBranchConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "branch", "v1"),
					resource.TestCheckResourceAttr("data.git_repository.test", "tag", "v1.0.0"),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_sha", "ea59fdb83b2ad18aa1da930f99aa34a1f325c009"),
				),
			},
		},
	})
}

func TestAccDataSourceRepository_AuthHTTP(t *testing.T) {
	repoURL := os.Getenv("GIT_REPO_URL")
	if os.Getenv("GIT_USERNAME") == "" || os.Getenv("GIT_PASSWORD") == "" || repoURL == "" {
		t.Skip(`Skipping test because "GIT_REPO_URL", "GIT_USERNAME" and "GIT_PASSWORD" is not set`)
	}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryURLConfig(repoURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.git_repository.test", "branch"),
					resource.TestCheckResourceAttrSet("data.git_repository.test", "commit_sha"),
				),
			},
		},
	})
}

func TestAccDataSourceRepository_AuthSSH(t *testing.T) {
	if os.Getenv("GIT_PRIVATE_KEY") == "" && os.Getenv("GIT_PRIVATE_KEY_FILE") == "" {
		t.Skip(`Skipping test because "GIT_PRIVATE_KEY" or "GIT_PRIVATE_KEY_FILE" is not set`)
	}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryURLConfig("git@github.com:innovationnorway/terraform-module-acctest.git"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "branch", "master"),
					resource.TestCheckResourceAttr("data.git_repository.test", "tag", "v0.1.0"),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_sha", "a48e01ced30c807293088746dc2b6715d78882f1"),
				),
			},
		},
	})
}

func testAccDataSourceRepositoryURLConfig(url string) string {
	return fmt.Sprintf(`
data "git_repository" "test" {
  url = "%s"
}
`, url)
}

func testAccDataSourceRepositoryPathConfig(path string) string {
	return fmt.Sprintf(`
data "git_repository" "test" {
  path = "%s"
}
`, path)
}

const testAccDataSourceRepositoryTagConfig = `
data "git_repository" "test" {
  url = "https://github.com/innovationnorway/terraform-module-acctest.git"
  tag = "v1.0.0"
}
`
const testAccDataSourceRepositoryBranchConfig = `
data "git_repository" "test" {
  url    = "https://github.com/innovationnorway/terraform-module-acctest.git"
  branch = "v1"
}
`
