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
)

//nolint:errcheck
func TestGenMarkdownTree(t *testing.T) {
	tmpdir := t.TempDir()

	err := copyFile(path.Join("fixtures", "buildx_stop.pre.md"), path.Join(tmpdir, "buildx_stop.md"))
	if err != nil {
		t.Fatal(err)
	}

	c, err := New(Options{
		Root:      buildxCmd,
		SourceDir: tmpdir,
		Plugin:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = c.GenMarkdownTree(buildxCmd)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range []string{"buildx.md", "buildx_build.md", "buildx_stop.md"} {
		tt := tt
		t.Run(tt, func(t *testing.T) {
			bres, err := os.ReadFile(filepath.Join(tmpdir, tt))
			if err != nil {
				t.Fatal(err)
			}

			bexc, err := os.ReadFile(path.Join("fixtures", tt))
			if err != nil {
				t.Fatal(err)
			}
			if string(bexc) != string(bres) {
				t.Fatalf("expected:\n%s\ngot:\n%s", string(bexc), string(bres))
			}
		})
	}
}
