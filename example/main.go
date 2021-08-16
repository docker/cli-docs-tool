// Copyright 2021 docgen authors
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

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/docker/buildx/commands"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docgen"
	"github.com/spf13/cobra"
)

const sourcePath = "docs/reference/"

func main() {
	log.SetFlags(0)

	dockerCLI, err := command.NewDockerCli()
	if err != nil {
		log.Printf("ERROR: %+v", err)
	}

	cmd := &cobra.Command{
		Use:               "docker [OPTIONS] COMMAND [ARG...]",
		Short:             "The base command for the Docker CLI.",
		DisableAutoGenTag: true,
	}

	cmd.AddCommand(commands.NewRootCmd("buildx", true, dockerCLI))
	docgen.DisableFlagsInUseLine(cmd)

	cwd, _ := os.Getwd()
	source := filepath.Join(cwd, sourcePath)

	if err = os.MkdirAll(source, 0755); err != nil {
		log.Printf("ERROR: %+v", err)
	}
	if err = docgen.GenTree(cmd, source); err != nil {
		log.Printf("ERROR: %+v", err)
	}
}
