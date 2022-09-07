terraform {
  required_providers {
    git = {
      source  = "Pango-inc/git"
      version = ">= 0.1.3"
    }
  }
}

data "git_repository" "example" {
}

output "repository" {
  value = {
    url           = data.git_repository.example.url
    branch        = data.git_repository.example.branch
    commit        = substr(data.git_repository.example.commit_sha, 0, 7)
    tag           = data.git_repository.example.tag
    relative_path = data.git_repository.example.relative_path
    clean         = data.git_repository.example.clean
  }
}
