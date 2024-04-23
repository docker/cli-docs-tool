// Copyright 2024 cli-docs-tool authors
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
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/cobra/doc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:errcheck
func TestGenManTree(t *testing.T) {
	setup()
	tmpdir := t.TempDir()

	epoch, err := time.Parse("2006-Jan-02", "2020-Jan-10")
	require.NoError(t, err)
	t.Setenv("SOURCE_DATE_EPOCH", strconv.FormatInt(epoch.Unix(), 10))

	require.NoError(t, copyFile(path.Join("fixtures", "buildx_stop.pre.md"), path.Join(tmpdir, "buildx_stop.md")))

	c, err := New(Options{
		Root:      dockerCmd,
		SourceDir: tmpdir,
		Plugin:    false,
		ManHeader: &doc.GenManHeader{
			Title:   "DOCKER",
			Section: "1",
			Source:  "Docker Community",
			Manual:  "Docker User Manuals",
		},
	})
	require.NoError(t, err)
	require.NoError(t, c.GenManTree(dockerCmd))

	seen := make(map[string]struct{})
	remanpage := regexp.MustCompile(`\.\d+$`)

	filepath.Walk("fixtures", func(path string, info fs.FileInfo, err error) error {
		fname := filepath.Base(path)
		// ignore dirs and any file that is not a manpage
		if info.IsDir() || !remanpage.MatchString(fname) {
			return nil
		}
		t.Run(fname, func(t *testing.T) {
			seen[fname] = struct{}{}
			require.NoError(t, err)

			bres, err := os.ReadFile(filepath.Join(tmpdir, fname))
			require.NoError(t, err)

			bexc, err := os.ReadFile(path)
			require.NoError(t, err)
			assert.Equal(t, string(bexc), string(bres))
		})
		return nil
	})

	filepath.Walk(tmpdir, func(path string, info fs.FileInfo, err error) error {
		fname := filepath.Base(path)
		// ignore dirs and any file that is not a manpage
		if info.IsDir() || !remanpage.MatchString(fname) {
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
