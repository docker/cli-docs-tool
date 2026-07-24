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
)

//nolint:errcheck
func TestGenManTree(t *testing.T) {
	setup()
	tmpdir := t.TempDir()

	epoch, err := time.Parse("2006-Jan-02", "2020-Jan-10")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("SOURCE_DATE_EPOCH", strconv.FormatInt(epoch.Unix(), 10))

	err = copyFile(path.Join("fixtures", "buildx_stop.pre.md"), path.Join(tmpdir, "buildx_stop.md"))
	if err != nil {
		t.Fatal(err)
	}

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
	if err != nil {
		t.Fatal(err)
	}

	err = c.GenManTree(dockerCmd)
	if err != nil {
		t.Fatal(err)
	}

	seen := make(map[string]struct{})
	remanpage := regexp.MustCompile(`\.\d+$`)

	_ = filepath.Walk("fixtures", func(path string, info fs.FileInfo, _ error) error {
		fname := filepath.Base(path)
		// ignore dirs and any file that is not a manpage
		if info.IsDir() || !remanpage.MatchString(fname) {
			return nil
		}
		t.Run(fname, func(t *testing.T) {
			seen[fname] = struct{}{}

			bres, err := os.ReadFile(filepath.Join(tmpdir, fname))
			if err != nil {
				t.Fatal(err)
			}

			bexc, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(bexc) != string(bres) {
				t.Fatalf("expected:\n%s\ngot:\n%s", string(bexc), string(bres))
			}
		})
		return nil
	})

	_ = filepath.Walk(tmpdir, func(path string, info fs.FileInfo, _ error) error {
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
