variable "GO_VERSION" {
  default = "1.16"
}

group "default" {
  targets = ["test"]
}

group "validate" {
  targets = ["lint", "vendor-validate"]
}

target "lint" {
  args = {
    GO_VERSION = GO_VERSION
  }
  dockerfile = "./hack/lint.Dockerfile"
  target = "lint"
  output = ["type=cacheonly"]
}

target "vendor-validate" {
  args = {
    GO_VERSION = GO_VERSION
  }
  dockerfile = "./hack/vendor.Dockerfile"
  target = "validate"
  output = ["type=cacheonly"]
}

target "vendor-update" {
  args = {
    GO_VERSION = GO_VERSION
  }
  dockerfile = "./hack/vendor.Dockerfile"
  target = "update"
  output = ["."]
}

target "test" {
  args = {
    GO_VERSION = GO_VERSION
  }
  dockerfile = "./hack/test.Dockerfile"
  target = "test-coverage"
  output = ["."]
}

variable "GITHUB_REPOSITORY" {}
variable "GITHUB_REF" {}
target "godev" {
  args = {
    GO_VERSION = GO_VERSION
    GITHUB_REPOSITORY = GITHUB_REPOSITORY
    GITHUB_REF = GITHUB_REF
  }
  dockerfile = "./hack/godev.Dockerfile"
  target = "godev"
  output = ["type=cacheonly"]
}
