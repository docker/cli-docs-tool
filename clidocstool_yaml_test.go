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
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:errcheck
func TestGenYamlTree(t *testing.T) {
	tmpdir := t.TempDir()

	c, err := New(Options{
		Root:      buildxCmd,
		SourceDir: tmpdir,
		Plugin:    true,
	})
	require.NoError(t, err)
	require.NoError(t, c.GenYamlTree(buildxCmd))

	for _, tt := range []string{"docker_buildx.yaml", "docker_buildx_build.yaml", "docker_buildx_install.yaml", "docker_buildx_stop.yaml"} {
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
