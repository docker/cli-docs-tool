// Copyright 2021 cli-docs-tool authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clidocstool

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/docker/cli-docs-tool/annotation"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dockerCmd      *cobra.Command
	buildxCmd      *cobra.Command
	buildxBuildCmd *cobra.Command
	buildxStopCmd  *cobra.Command
)

//nolint:errcheck
func init() {
	dockerCmd = &cobra.Command{
		Use:                   "docker [OPTIONS] COMMAND [ARG...]",
		Short:                 "A self-sufficient runtime for containers",
		SilenceUsage:          true,
		SilenceErrors:         true,
		TraverseChildren:      true,
		Run:                   func(cmd *cobra.Command, args []string) {},
		Version:               "20.10.8",
		DisableFlagsInUseLine: true,
	}
	buildxCmd = &cobra.Command{
		Use:   "buildx",
		Short: "Docker Buildx",
		Long:  `Extended build capabilities with BuildKit`,
		Annotations: map[string]string{
			annotation.CodeDelimiter: `"`,
		},
	}
	buildxBuildCmd = &cobra.Command{
		Use:     "build [OPTIONS] PATH | URL | -",
		Aliases: []string{"b"},
		Short:   "Start a build",
		Run:     func(cmd *cobra.Command, args []string) {},
	}
	buildxStopCmd = &cobra.Command{
		Use:   "stop [NAME]",
		Short: "Stop builder instance",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	buildxPFlags := buildxCmd.PersistentFlags()
	buildxPFlags.String("builder", os.Getenv("BUILDX_BUILDER"), "Override the configured builder instance")

	buildxBuildFlags := buildxBuildCmd.Flags()

	var ignore string
	var ignoreSlice []string
	var ignoreBool bool
	var ignoreInt int64

	buildxBuildFlags.StringSlice("add-host", []string{}, `Add a custom host-to-IP mapping (format: 'host:ip')`)
	buildxBuildFlags.SetAnnotation("add-host", annotation.ExternalURL, []string{"https://docs.docker.com/engine/reference/commandline/build/#add-entries-to-container-hosts-file---add-host"})
	buildxBuildFlags.SetAnnotation("add-host", annotation.CodeDelimiter, []string{`'`})

	buildxBuildFlags.StringSlice("allow", []string{}, `Allow extra privileged entitlement (e.g., "network.host", "security.insecure")`)

	buildxBuildFlags.StringArray("build-arg", []string{}, "Set build-time variables")
	buildxBuildFlags.SetAnnotation("build-arg", annotation.ExternalURL, []string{"https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg"})

	buildxBuildFlags.StringArray("cache-from", []string{}, `External cache sources (e.g., "user/app:cache", "type=local,src=path/to/dir")`)

	buildxBuildFlags.StringArray("cache-to", []string{}, `Cache export destinations (e.g., "user/app:cache", "type=local,dest=path/to/dir")`)

	buildxBuildFlags.String("cgroup-parent", "", "Optional parent cgroup for the container")
	buildxBuildFlags.SetAnnotation("cgroup-parent", annotation.ExternalURL, []string{"https://docs.docker.com/engine/reference/commandline/build/#use-a-custom-parent-cgroup---cgroup-parent"})

	buildxBuildFlags.StringP("file", "f", "", `Name of the Dockerfile (default: "PATH/Dockerfile")`)
	buildxBuildFlags.SetAnnotation("file", annotation.ExternalURL, []string{"https://docs.docker.com/engine/reference/commandline/build/#specify-a-dockerfile--f"})

	buildxBuildFlags.String("iidfile", "", "Write the image ID to the file")

	buildxBuildFlags.StringArray("label", []string{}, "Set metadata for an image")

	buildxBuildFlags.Bool("load", false, `Shorthand for "--output=type=docker"`)

	buildxBuildFlags.String("network", "default", `Set the networking mode for the "RUN" instructions during build`)

	buildxBuildFlags.StringArrayP("output", "o", []string{}, `Output destination (format: "type=local,dest=path")`)

	buildxBuildFlags.StringArray("platform", []string{}, "Set target platform for build")

	buildxBuildFlags.Bool("push", false, `Shorthand for "--output=type=registry"`)

	buildxBuildFlags.BoolP("quiet", "q", false, "Suppress the build output and print image ID on success")

	buildxBuildFlags.StringArray("secret", []string{}, `Secret file to expose to the build (format: "id=mysecret,src=/local/secret")`)

	buildxBuildFlags.StringVar(&ignore, "shm-size", "", `Size of "/dev/shm"`)

	buildxBuildFlags.StringArray("ssh", []string{}, `SSH agent socket or keys to expose to the build (format: "default|<id>[=<socket>|<key>[,<key>]]")`)

	buildxBuildFlags.StringArrayP("tag", "t", []string{}, `Name and optionally a tag (format: "name:tag")`)
	buildxBuildFlags.SetAnnotation("tag", annotation.ExternalURL, []string{"https://docs.docker.com/engine/reference/commandline/build/#tag-an-image--t"})

	buildxBuildFlags.String("target", "", "Set the target build stage to build.")
	buildxBuildFlags.SetAnnotation("target", annotation.ExternalURL, []string{"https://docs.docker.com/engine/reference/commandline/build/#specifying-target-build-stage---target"})

	buildxBuildFlags.StringVar(&ignore, "ulimit", "", "Ulimit options")

	// hidden flags
	buildxBuildFlags.BoolVar(&ignoreBool, "compress", false, "Compress the build context using gzip")
	buildxBuildFlags.MarkHidden("compress")

	buildxBuildFlags.StringVar(&ignore, "isolation", "", "Container isolation technology")
	buildxBuildFlags.MarkHidden("isolation")
	buildxBuildFlags.SetAnnotation("isolation", "flag-warn", []string{"isolation flag is deprecated with BuildKit."})

	buildxBuildFlags.StringSliceVar(&ignoreSlice, "security-opt", []string{}, "Security options")
	buildxBuildFlags.MarkHidden("security-opt")
	buildxBuildFlags.SetAnnotation("security-opt", "flag-warn", []string{`security-opt flag is deprecated. "RUN --security=insecure" should be used with BuildKit.`})

	buildxBuildFlags.BoolVar(&ignoreBool, "squash", false, "Squash newly built layers into a single new layer")
	buildxBuildFlags.MarkHidden("squash")
	buildxBuildFlags.SetAnnotation("squash", "flag-warn", []string{"experimental flag squash is removed with BuildKit. You should squash inside build using a multi-stage Dockerfile for efficiency."})

	buildxBuildFlags.StringVarP(&ignore, "memory", "m", "", "Memory limit")
	buildxBuildFlags.MarkHidden("memory")

	buildxBuildFlags.StringVar(&ignore, "memory-swap", "", `Swap limit equal to memory plus swap: "-1" to enable unlimited swap`)
	buildxBuildFlags.MarkHidden("memory-swap")

	buildxBuildFlags.Int64VarP(&ignoreInt, "cpu-shares", "c", 0, "CPU shares (relative weight)")
	buildxBuildFlags.MarkHidden("cpu-shares")

	buildxBuildFlags.Int64Var(&ignoreInt, "cpu-period", 0, "Limit the CPU CFS (Completely Fair Scheduler) period")
	buildxBuildFlags.MarkHidden("cpu-period")

	buildxBuildFlags.Int64Var(&ignoreInt, "cpu-quota", 0, "Limit the CPU CFS (Completely Fair Scheduler) quota")
	buildxBuildFlags.MarkHidden("cpu-quota")

	buildxBuildFlags.StringVar(&ignore, "cpuset-cpus", "", `CPUs in which to allow execution ("0-3", "0,1")`)
	buildxBuildFlags.MarkHidden("cpuset-cpus")

	buildxBuildFlags.StringVar(&ignore, "cpuset-mems", "", `MEMs in which to allow execution ("0-3", "0,1")`)
	buildxBuildFlags.MarkHidden("cpuset-mems")

	buildxBuildFlags.BoolVar(&ignoreBool, "rm", true, "Remove intermediate containers after a successful build")
	buildxBuildFlags.MarkHidden("rm")

	buildxBuildFlags.BoolVar(&ignoreBool, "force-rm", false, "Always remove intermediate containers")
	buildxBuildFlags.MarkHidden("force-rm")

	buildxCmd.AddCommand(buildxBuildCmd)
	buildxCmd.AddCommand(buildxStopCmd)
	dockerCmd.AddCommand(buildxCmd)
}

//nolint:errcheck
func TestGenAllTree(t *testing.T) {
	tmpdir := t.TempDir()

	err := copyFile(path.Join("fixtures", "buildx_stop.pre.md"), path.Join(tmpdir, "buildx_stop.md"))
	require.NoError(t, err)

	c, err := New(Options{
		Root:      buildxCmd,
		SourceDir: tmpdir,
		Plugin:    true,
	})
	require.NoError(t, err)
	require.NoError(t, c.GenAllTree())

	for _, tt := range []string{"buildx.md", "buildx_build.md", "buildx_stop.md", "docker_buildx.yaml", "docker_buildx_build.yaml", "docker_buildx_stop.yaml"} {
		tt := tt
		t.Run(tt, func(t *testing.T) {
			bres, err := os.ReadFile(filepath.Join(tmpdir, tt))
			require.NoError(t, err)

			bexc, err := os.ReadFile(path.Join("fixtures", tt))
			require.NoError(t, err)
			assert.Equal(t, string(bexc), string(bres))
		})
	}
}
