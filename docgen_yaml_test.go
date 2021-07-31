package docgen

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
