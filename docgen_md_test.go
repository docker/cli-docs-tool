package docgen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestGenMarkdownTree(t *testing.T) {
	c := &cobra.Command{Use: "do [OPTIONS] arg1 arg2"}
	s := &cobra.Command{Use: "sub [OPTIONS] arg1 arg2"}
	c.AddCommand(s)

	tmpdir, err := ioutil.TempDir("", "test-gen-markdown-tree")
	if err != nil {
		t.Fatalf("Failed to create tmpdir: %s", err.Error())
	}
	defer os.RemoveAll(tmpdir)

	if err := GenMarkdownTree(c, tmpdir); err != nil {
		t.Fatalf("GenMarkdownTree failed: %s", err.Error())
	}

	if _, err := os.Stat(filepath.Join(tmpdir, "sub.md")); err != nil {
		t.Fatalf("Expected file 'sub.md' to exist")
	}
}
