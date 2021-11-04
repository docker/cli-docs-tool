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

	flags := buildxBuildCmd.Flags()
	flags.Bool("push", false, "Shorthand for --output=type=registry")
	flags.Bool("load", false, "Shorthand for --output=type=docker")
	flags.StringArrayP("tag", "t", []string{}, "Name and optionally a tag in the 'name:tag' format")
	flags.SetAnnotation("tag", "docs.external.url", []string{"https://docs.docker.com/engine/reference/commandline/build/#tag-an-image--t"})
	flags.StringArray("build-arg", []string{}, "Set build-time variables")
	flags.SetAnnotation("build-arg", "docs.external.url", []string{"https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg"})
	flags.StringP("file", "f", "", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")
	flags.SetAnnotation("file", "docs.external.url", []string{"https://docs.docker.com/engine/reference/commandline/build/#specify-a-dockerfile--f"})
	flags.StringArray("label", []string{}, "Set metadata for an image")
	flags.StringArray("cache-from", []string{}, "External cache sources (eg. user/app:cache, type=local,src=path/to/dir)")
	flags.StringArray("cache-to", []string{}, "Cache export destinations (eg. user/app:cache, type=local,dest=path/to/dir)")
	flags.String("target", "", "Set the target build stage to build.")
	flags.SetAnnotation("target", "docs.external.url", []string{"https://docs.docker.com/engine/reference/commandline/build/#specifying-target-build-stage---target"})
	flags.StringSlice("allow", []string{}, "Allow extra privileged entitlement, e.g. network.host, security.insecure")
	flags.StringArray("platform", []string{}, "Set target platform for build")
	flags.StringArray("secret", []string{}, "Secret file to expose to the build: id=mysecret,src=/local/secret")
	flags.StringArray("ssh", []string{}, "SSH agent socket or keys to expose to the build (format: `default|<id>[=<socket>|<key>[,<key>]]`)")
	flags.StringArrayP("output", "o", []string{}, "Output destination (format: type=local,dest=path)")
	// not implemented
	flags.String("network", "default", "Set the networking mode for the RUN instructions during build")
	flags.StringSlice("add-host", []string{}, "Add a custom host-to-IP mapping (host:ip)")
	flags.SetAnnotation("add-host", "docs.external.url", []string{"https://docs.docker.com/engine/reference/commandline/build/#add-entries-to-container-hosts-file---add-host"})
	flags.String("iidfile", "", "Write the image ID to the file")
	// hidden flags
	flags.BoolP("quiet", "q", false, "Suppress the build output and print image ID on success")
	flags.MarkHidden("quiet")
	flags.Bool("squash", false, "Squash newly built layers into a single new layer")
	flags.MarkHidden("squash")
	flags.String("ulimit", "", "Ulimit options")
	flags.MarkHidden("ulimit")
	flags.StringSlice("security-opt", []string{}, "Security options")
	flags.MarkHidden("security-opt")
	flags.Bool("compress", false, "Compress the build context using gzip")

	buildxCmd.AddCommand(buildxBuildCmd)
	dockerCmd.AddCommand(buildxCmd)
}

//nolint:errcheck
func TestGenAllTree(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test-gen-all-tree")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	c, err := New(Options{
		Root:      buildxCmd,
		SourceDir: tmpdir,
		Plugin:    true,
	})
	require.NoError(t, err)
	require.NoError(t, c.GenAllTree())

	for _, tt := range []string{"buildx.md", "buildx_build.md", "docker_buildx.yaml", "docker_buildx_build.yaml"} {
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
