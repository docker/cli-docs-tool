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
)

func TestGenYamlTree(t *testing.T) {
	c := &cobra.Command{Use: "do [OPTIONS] arg1 arg2"}
	s := &cobra.Command{Use: "sub [OPTIONS] arg1 arg2", Run: func(cmd *cobra.Command, args []string) {}}
	c.AddCommand(s)

	tmpdir, err := ioutil.TempDir("", "test-gen-yaml-tree")
	if err != nil {
		t.Fatalf("Failed to create tmpdir: %s", err.Error())
	}
	defer os.RemoveAll(tmpdir)

	if err := GenYamlTree(c, tmpdir); err != nil {
		t.Fatalf("GenYamlTree failed: %s", err.Error())
	}

	if _, err := os.Stat(filepath.Join(tmpdir, "do_sub.yaml")); err != nil {
		t.Fatalf("Expected file 'do_sub.yaml' to exist")
	}
}