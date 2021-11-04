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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

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
		Short: "Build with BuildKit",
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
	buildxBuildFlags.Bool("push", false, "Shorthand for --output=type=registry")
	buildxBuildFlags.Bool("load", false, "Shorthand for --output=type=docker")
	buildxBuildFlags.StringArrayP("tag", "t", []string{}, "Name and optionally a tag in the 'name:tag' format")
	buildxBuildFlags.SetAnnotation("tag", AnnotationExternalUrl, []string{"https://docs.docker.com/engine/reference/commandline/build/#tag-an-image--t"})
	buildxBuildFlags.StringArray("build-arg", []string{}, "Set build-time variables")
	buildxBuildFlags.SetAnnotation("build-arg", AnnotationExternalUrl, []string{"https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg"})
	buildxBuildFlags.StringP("file", "f", "", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")
	buildxBuildFlags.SetAnnotation("file", AnnotationExternalUrl, []string{"https://docs.docker.com/engine/reference/commandline/build/#specify-a-dockerfile--f"})
	buildxBuildFlags.StringArray("label", []string{}, "Set metadata for an image")
	buildxBuildFlags.StringArray("cache-from", []string{}, "External cache sources (eg. user/app:cache, type=local,src=path/to/dir)")
	buildxBuildFlags.StringArray("cache-to", []string{}, "Cache export destinations (eg. user/app:cache, type=local,dest=path/to/dir)")
	buildxBuildFlags.String("target", "", "Set the target build stage to build.")
	buildxBuildFlags.SetAnnotation("target", AnnotationExternalUrl, []string{"https://docs.docker.com/engine/reference/commandline/build/#specifying-target-build-stage---target"})
	buildxBuildFlags.StringSlice("allow", []string{}, "Allow extra privileged entitlement, e.g. network.host, security.insecure")
	buildxBuildFlags.StringArray("platform", []string{}, "Set target platform for build")
	buildxBuildFlags.StringArray("secret", []string{}, "Secret file to expose to the build: id=mysecret,src=/local/secret")
	buildxBuildFlags.StringArray("ssh", []string{}, "SSH agent socket or keys to expose to the build (format: `default|<id>[=<socket>|<key>[,<key>]]`)")
	buildxBuildFlags.StringArrayP("output", "o", []string{}, "Output destination (format: type=local,dest=path)")
	// not implemented
	buildxBuildFlags.String("network", "default", "Set the networking mode for the RUN instructions during build")
	buildxBuildFlags.StringSlice("add-host", []string{}, "Add a custom host-to-IP mapping (host:ip)")
	buildxBuildFlags.SetAnnotation("add-host", AnnotationExternalUrl, []string{"https://docs.docker.com/engine/reference/commandline/build/#add-entries-to-container-hosts-file---add-host"})
	buildxBuildFlags.String("iidfile", "", "Write the image ID to the file")
	// hidden flags
	buildxBuildFlags.BoolP("quiet", "q", false, "Suppress the build output and print image ID on success")
	buildxBuildFlags.MarkHidden("quiet")
	buildxBuildFlags.Bool("squash", false, "Squash newly built layers into a single new layer")
	buildxBuildFlags.MarkHidden("squash")
	buildxBuildFlags.String("ulimit", "", "Ulimit options")
	buildxBuildFlags.MarkHidden("ulimit")
	buildxBuildFlags.StringSlice("security-opt", []string{}, "Security options")
	buildxBuildFlags.MarkHidden("security-opt")
	buildxBuildFlags.Bool("compress", false, "Compress the build context using gzip")

	buildxCmd.AddCommand(buildxBuildCmd)
	buildxCmd.AddCommand(buildxStopCmd)
	dockerCmd.AddCommand(buildxCmd)
}

//nolint:errcheck
func TestGenAllTree(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test-gen-all-tree")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	err = copyFile(path.Join("fixtures", "buildx_stop.pre.md"), path.Join(tmpdir, "buildx_stop.md"))
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
			fres := filepath.Join(tmpdir, tt)
			require.FileExists(t, fres)
			bres, err := ioutil.ReadFile(fres)
			require.NoError(t, err)

			bexc, err := ioutil.ReadFile(path.Join("fixtures", tt))
			require.NoError(t, err)
			assert.Equal(t, string(bexc), string(bres))
		})
	}
}
