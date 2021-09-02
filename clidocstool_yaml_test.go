// Copyright 2017 cli-docs-tool authors
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
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:errcheck
func TestGenYamlTree(t *testing.T) {
	c := &cobra.Command{Use: "do [OPTIONS] arg1 arg2"}
	s := &cobra.Command{Use: "sub [OPTIONS] arg1 arg2", Run: func(cmd *cobra.Command, args []string) {}}

	flags := s.Flags()
	_ = flags.Bool("push", false, "Shorthand for --output=type=registry")
	_ = flags.Bool("load", false, "Shorthand for --output=type=docker")
	_ = flags.StringArrayP("tag", "t", []string{}, "Name and optionally a tag in the 'name:tag' format")
	flags.SetAnnotation("tag", "docs.external.url", []string{"https://docs.docker.com/engine/reference/commandline/build/#tag-an-image--t"})
	_ = flags.StringArray("build-arg", []string{}, "Set build-time variables")
	flags.SetAnnotation("build-arg", "docs.external.url", []string{"https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg"})
	_ = flags.StringP("file", "f", "", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")
	flags.SetAnnotation("file", "docs.external.url", []string{"https://docs.docker.com/engine/reference/commandline/build/#specify-a-dockerfile--f"})
	_ = flags.StringArray("label", []string{}, "Set metadata for an image")
	_ = flags.StringArray("cache-from", []string{}, "External cache sources (eg. user/app:cache, type=local,src=path/to/dir)")
	_ = flags.StringArray("cache-to", []string{}, "Cache export destinations (eg. user/app:cache, type=local,dest=path/to/dir)")
	_ = flags.String("target", "", "Set the target build stage to build.")
	flags.SetAnnotation("target", "docs.external.url", []string{"https://docs.docker.com/engine/reference/commandline/build/#specifying-target-build-stage---target"})
	_ = flags.StringSlice("allow", []string{}, "Allow extra privileged entitlement, e.g. network.host, security.insecure")
	_ = flags.StringArray("platform", []string{}, "Set target platform for build")
	_ = flags.StringArray("secret", []string{}, "Secret file to expose to the build: id=mysecret,src=/local/secret")
	_ = flags.StringArray("ssh", []string{}, "SSH agent socket or keys to expose to the build (format: `default|<id>[=<socket>|<key>[,<key>]]`)")
	_ = flags.StringArrayP("output", "o", []string{}, "Output destination (format: type=local,dest=path)")
	// not implemented
	_ = flags.String("network", "default", "Set the networking mode for the RUN instructions during build")
	_ = flags.StringSlice("add-host", []string{}, "Add a custom host-to-IP mapping (host:ip)")
	_ = flags.SetAnnotation("add-host", "docs.external.url", []string{"https://docs.docker.com/engine/reference/commandline/build/#add-entries-to-container-hosts-file---add-host"})
	_ = flags.String("iidfile", "", "Write the image ID to the file")
	// hidden flags
	_ = flags.BoolP("quiet", "q", false, "Suppress the build output and print image ID on success")
	flags.MarkHidden("quiet")
	_ = flags.Bool("squash", false, "Squash newly built layers into a single new layer")
	flags.MarkHidden("squash")
	_ = flags.String("ulimit", "", "Ulimit options")
	flags.MarkHidden("ulimit")
	_ = flags.StringSlice("security-opt", []string{}, "Security options")
	flags.MarkHidden("security-opt")
	_ = flags.Bool("compress", false, "Compress the build context using gzip")

	c.AddCommand(s)

	tmpdir, err := ioutil.TempDir("", "test-gen-yaml-tree")
	require.NoError(t, err)

	defer os.RemoveAll(tmpdir)
	require.NoError(t, GenYamlTree(c, tmpdir))

	fres := filepath.Join(tmpdir, "do_sub.yaml")
	require.FileExists(t, fres)
	bres, err := ioutil.ReadFile(fres)
	require.NoError(t, err)
	bexc, err := ioutil.ReadFile("fixtures/do_sub.yaml")
	require.NoError(t, err)
	assert.Equal(t, string(bres), string(bexc))
}
