package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/skooma-cli/skooma/internal/templates"
	"github.com/spf13/cobra"
)

var templateLsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List available templates",
	Long:    `List all available templates that can be used with the brew command.`,
	Aliases: []string{"list"},
	Run: func(cmd *cobra.Command, args []string) {
		templates, err := templates.GetTemplates()
		if err != nil {
			fmt.Printf("Error loading templates: %v\n", err)
			return
		}

		if len(templates) == 0 {
			fmt.Println("No templates available. Use 'skooma template add' to add a template.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tDESCRIPTION\tREPO\tAUTHOR")

		for name, tmpl := range templates {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, tmpl.Description, tmpl.RepoURL, tmpl.Author)
		}

		w.Flush()
	},
}
