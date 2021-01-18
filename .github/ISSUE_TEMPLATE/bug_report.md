---
name: Bug report
about: Create a report to help us improve
title: ''
labels: bug
assignees: ''

---

Thank you 🙇‍♀ for opening an issue. If your issue is relevant to this repository, please include the information below:

### Terraform Version
Run `terraform -v` to show the version. If you are not running the latest version of Terraform, please upgrade because your issue may have already been fixed.

### Affected Resource(s)
Please list the resources as a list, for example:
- git_repository
- git_repository_file

If this issue appears to affect multiple resources, it may be an issue with Terraform's core, so please mention this.

### Terraform Configuration Files
```hcl
# Copy-paste your Terraform configurations here.
# For large Terraform configs, please provide a link to a GitHub Gist.
```

### Debug Output
Please provide a link to a [GitHub Gist](https://gist.github.com/) containing the complete debug output: https://www.terraform.io/docs/internals/debugging.html. Please do NOT paste the debug output in the issue; just paste a link to the Gist.

### Panic Output
If Terraform produced a panic, please provide a link to a [GitHub Gist](https://gist.github.com/) containing the output of the `crash.log`.

### Expected Behavior
What should have happened?

### Actual Behavior
What actually happened?

### Steps to Reproduce
Please list the steps required to reproduce the issue, for example:
1. `terraform apply`

### References
Are there any other GitHub issues (open or closed) or Pull Requests that should be linked here? For example:
- #1234
