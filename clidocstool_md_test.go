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
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:errcheck
func TestGenMarkdownTree(t *testing.T) {
	tmpdir := t.TempDir()

	require.NoError(t, copyFile(path.Join("fixtures", "buildx_stop.pre.md"), path.Join(tmpdir, "buildx_stop.md")))

	c, err := New(Options{
		Root:      buildxCmd,
		SourceDir: tmpdir,
		Plugin:    true,
	})
	require.NoError(t, err)
	require.NoError(t, c.GenMarkdownTree(buildxCmd))

	seen := make(map[string]struct{})

	filepath.Walk("fixtures", func(path string, info fs.FileInfo, err error) error {
		fname := filepath.Base(path)
		// ignore dirs, .pre.md files and any file that is not a .md file
		if info.IsDir() || !strings.HasSuffix(fname, ".md") || strings.HasSuffix(fname, ".pre.md") {
			return nil
		}
		t.Run(fname, func(t *testing.T) {
			seen[fname] = struct{}{}
			require.NoError(t, err)

			bres, err := os.ReadFile(filepath.Join(tmpdir, fname))
			require.NoError(t, err)

			bexc, err := os.ReadFile(path)
			require.NoError(t, err)
			require.Equal(t, string(bexc), string(bres))
		})
		return nil
	})

	filepath.Walk(tmpdir, func(path string, info fs.FileInfo, err error) error {
		fname := filepath.Base(path)
		// ignore dirs, .pre.md files and any file that is not a .md file
		if info.IsDir() || !strings.HasSuffix(fname, ".md") || strings.HasSuffix(fname, ".pre.md") {
			return nil
		}
		t.Run("seen_"+fname, func(t *testing.T) {
			if _, ok := seen[fname]; !ok {
				t.Errorf("file %s not found in fixtures", fname)
			}
		})
		return nil
	})
}
