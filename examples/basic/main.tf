terraform {
  required_providers {
    git = {
      source  = "innovationnorway/git"
      version = ">= 0.1.3"
    }
  }
}

data "git_repository" "example" {
  path = "../../"
}

output "repository" {
  value = {
    branch = data.git_repository.example.branch
    commit = substr(data.git_repository.example.commit_sha, 0, 7)
    tag    = data.git_repository.example.tag
  }
}
