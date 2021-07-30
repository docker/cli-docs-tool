package docgen

import (
	"github.com/spf13/cobra"
)

func GenTree(cmd *cobra.Command, dir string) error {
	var err error
	if err = GenMarkdownTree(cmd, dir); err != nil {
		return err
	}
	if err = GenYamlTree(cmd, dir); err != nil {
		return err
	}
	return nil
}
