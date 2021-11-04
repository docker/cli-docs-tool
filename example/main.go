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

package main

import (
	"log"
	"os"

	"github.com/docker/buildx/commands"
	clidocstool "github.com/docker/cli-docs-tool"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	pluginName        = "buildx"
	defaultSourcePath = "docs/"
)

type options struct {
	source string
	target string
}

func gen(opts *options) error {
	log.SetFlags(0)

	// create a new instance of Docker CLI
	dockerCLI, err := command.NewDockerCli()
	if err != nil {
		return err
	}

	// root command
	cmd := &cobra.Command{
		Use:   "buildx",
		Short: "Build with BuildKit",
	}

	// subcommand for the plugin
	cmd.AddCommand(commands.NewRootCmd(pluginName, true, dockerCLI))

	// create a new instance of cli-docs-tool
	c, err := clidocstool.New(clidocstool.Options{
		Root:      cmd,
		SourceDir: opts.source,
		TargetDir: opts.target,
		Plugin:    true,
	})
	if err != nil {
		return err
	}

	// generate all supported docs formats
	return c.GenAllTree()
}

func run() error {
	opts := &options{}
	flags := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	flags.StringVar(&opts.source, "source", defaultSourcePath, "Docs source folder")
	flags.StringVar(&opts.target, "target", defaultSourcePath, "Docs target folder")
	if err := flags.Parse(os.Args[1:]); err != nil {
		return err
	}
	return gen(opts)
}

func main() {
	if err := run(); err != nil {
		log.Printf("ERROR: %+v", err)
		os.Exit(1)
	}
}
