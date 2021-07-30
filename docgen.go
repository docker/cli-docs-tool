package docgen

import (
	"github.com/spf13/cobra"
)

// GenTree creates yaml and markdown structured ref files for this command
// and all descendants in the directory given. This function will just
// call GenMarkdownTree and GenYamlTree functions successively.
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

// VisitAll will traverse all commands from the root.
// This is different from the VisitAll of cobra.Command where only parents
// are checked.
func VisitAll(root *cobra.Command, fn func(*cobra.Command)) {
	for _, cmd := range root.Commands() {
		VisitAll(cmd, fn)
	}
	fn(root)
}

// DisableFlagsInUseLine sets the DisableFlagsInUseLine flag on all
// commands within the tree rooted at cmd.
func DisableFlagsInUseLine(cmd *cobra.Command) {
	VisitAll(cmd, func(ccmd *cobra.Command) {
		// do not add a `[flags]` to the end of the usage line.
		ccmd.DisableFlagsInUseLine = true
	})
}
