module github.com/crazy-max/yamldocs/example

go 1.16

require (
	github.com/crazy-max/yamldocs v0.0.0
	github.com/docker/buildx v0.6.0
	github.com/docker/cli v20.10.7+incompatible
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
)

replace (
	github.com/crazy-max/yamldocs => ../
	github.com/docker/cli => github.com/docker/cli v20.10.3-0.20210702143511-f782d1355eff+incompatible
	github.com/docker/docker => github.com/docker/docker v20.10.3-0.20210609100121-ef4d47340142+incompatible
)
