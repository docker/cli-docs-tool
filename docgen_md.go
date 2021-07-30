package docgen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GenMarkdown(cmd *cobra.Command, dir string) error {
	for _, c := range cmd.Commands() {
		if err := GenMarkdown(c, dir); err != nil {
			return err
		}
	}
	if !cmd.HasParent() {
		return nil
	}

	mdFile := mdFilename(cmd)
	fullPath := filepath.Join(dir, mdFile)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("INFO: Init markdown for %s", cmd.CommandPath())
		var icBuf bytes.Buffer
		icTpl, err := template.New("ic").Option("missingkey=error").Parse(`# {{ .Command }}

<!---MARKER_GEN_START-->
<!---MARKER_GEN_END-->

`)
		if err != nil {
			return err
		}
		if err = icTpl.Execute(&icBuf, struct {
			Command string
		}{
			Command: cmd.CommandPath(),
		}); err != nil {
			return err
		}
		if err = ioutil.WriteFile(fullPath, icBuf.Bytes(), 0644); err != nil {
			return err
		}
	}

	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}

	cs := string(content)

	start := strings.Index(cs, "<!---MARKER_GEN_START-->")
	end := strings.Index(cs, "<!---MARKER_GEN_END-->")

	if start == -1 {
		return errors.Errorf("no start marker in %s", mdFile)
	}
	if end == -1 {
		return errors.Errorf("no end marker in %s", mdFile)
	}

	out, err := mdCmdOutput(cmd, cs)
	if err != nil {
		return err
	}
	cont := cs[:start] + "<!---MARKER_GEN_START-->" + "\n" + out + "\n" + cs[end:]

	fi, err := os.Stat(fullPath)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fullPath, []byte(cont), fi.Mode()); err != nil {
		return errors.Wrapf(err, "failed to write %s", fullPath)
	}

	log.Printf("INFO: Markdown updated for %s", cmd.CommandPath())
	return nil
}

func mdFilename(cmd *cobra.Command) string {
	name := cmd.CommandPath()
	if i := strings.Index(name, " "); i >= 0 {
		name = name[i+1:]
	}
	return strings.ReplaceAll(name, " ", "_") + ".md"
}

func mdMakeLink(txt, link string) string {
	return "[" + txt + "](#" + link + ")"
}

func mdCmdOutput(cmd *cobra.Command, old string) (string, error) {
	b := &strings.Builder{}

	desc := cmd.Short
	if cmd.Long != "" {
		desc = cmd.Long
	}
	if desc != "" {
		fmt.Fprintf(b, "%s\n\n", desc)
	}

	if len(cmd.Aliases) != 0 {
		fmt.Fprintf(b, "### Aliases\n\n`%s`", cmd.Name())
		for _, a := range cmd.Aliases {
			fmt.Fprintf(b, ", `%s`", a)
		}
		fmt.Fprint(b, "\n\n")
	}

	if len(cmd.Commands()) != 0 {
		fmt.Fprint(b, "### Subcommands\n\n")
		fmt.Fprint(b, "| Name | Description |\n")
		fmt.Fprint(b, "| --- | --- |\n")
		for _, c := range cmd.Commands() {
			fmt.Fprintf(b, "| [`%s`](%s) | %s |\n", c.Name(), mdFilename(c), c.Short)
		}
		fmt.Fprint(b, "\n\n")
	}

	hasFlags := cmd.Flags().HasAvailableFlags()

	cmd.Flags().AddFlagSet(cmd.InheritedFlags())

	if hasFlags {
		fmt.Fprint(b, "### Options\n\n")
		fmt.Fprint(b, "| Name | Description |\n")
		fmt.Fprint(b, "| --- | --- |\n")

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Hidden {
				return
			}
			isLink := strings.Contains(old, "<a name=\""+f.Name+"\"></a>")
			fmt.Fprint(b, "| ")
			if f.Shorthand != "" {
				name := "`-" + f.Shorthand + "`"
				if isLink {
					name = mdMakeLink(name, f.Name)
				}
				fmt.Fprintf(b, "%s, ", name)
			}
			name := "`--" + f.Name
			if f.Value.Type() != "bool" {
				name += " " + f.Value.Type()
			}
			name += "`"
			if isLink {
				name = mdMakeLink(name, f.Name)
			}
			fmt.Fprintf(b, "%s | %s |\n", name, f.Usage)
		})
		fmt.Fprintln(b, "")
	}

	return b.String(), nil
}
