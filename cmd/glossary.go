package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/clioutput"
	"github.com/joshk418/flexspec/internal/glossary"
)

var (
	glossaryJSON  bool
	addDefinition string
	addAliases    []string
	addCategory   string
	addSources    []string
)

var glossaryCmd = &cobra.Command{
	Use:   "glossary",
	Short: "List, query, and manage project glossary terms",
	Long: `List, query, and add entries to .flexspec/glossary.yaml.

Use --json for machine-readable output suitable for scripts and agents.`,
}

var glossaryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all glossary terms",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("resolve working directory: %w", err)
		}

		doc, err := glossary.Load(root)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if glossaryJSON {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(doc)
		}

		if len(doc.Terms) == 0 {
			_, _ = fmt.Fprintln(out, "No glossary terms")
			return nil
		}

		rows := make([][]string, 0, len(doc.Terms))
		for _, e := range doc.Terms {
			rows = append(rows, []string{
				e.Term,
				e.Definition,
				displayOrDash(e.Category),
			})
		}
		return clioutput.WriteTable(out,
			[]string{"TERM", "DEFINITION", "CATEGORY"},
			rows,
		)
	},
}

var glossaryQueryCmd = &cobra.Command{
	Use:   "query <text>",
	Short: "Search glossary terms",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("query text is required")
		}
		root, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("resolve working directory: %w", err)
		}

		doc, err := glossary.Load(root)
		if err != nil {
			return err
		}

		results := glossary.Query(doc, args[0])
		out := cmd.OutOrStdout()
		if glossaryJSON {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(results)
		}

		if len(results) == 0 {
			_, _ = fmt.Fprintln(out, "No matches")
			return nil
		}

		rows := make([][]string, 0, len(results))
		for _, e := range results {
			rows = append(rows, []string{
				e.Term,
				e.Definition,
				displayOrDash(e.Category),
			})
		}
		return clioutput.WriteTable(out,
			[]string{"TERM", "DEFINITION", "CATEGORY"},
			rows,
		)
	},
}

var glossaryAddCmd = &cobra.Command{
	Use:   "add <term>",
	Short: "Add or update a glossary term",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("term is required")
		}
		root, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("resolve working directory: %w", err)
		}

		term := strings.TrimSpace(args[0])
		if term == "" {
			return fmt.Errorf("term is required")
		}
		if strings.TrimSpace(addDefinition) == "" {
			return fmt.Errorf("--definition is required")
		}

		doc, err := glossary.Load(root)
		if err != nil {
			return err
		}

		entry := glossary.Entry{
			Term:       term,
			Definition: addDefinition,
			Category:   addCategory,
			Aliases:    addAliases,
			Sources:    addSources,
		}
		doc, err = glossary.Upsert(doc, entry)
		if err != nil {
			return err
		}
		if err := glossary.Save(root, doc); err != nil {
			return err
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Added %q\n", term)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(glossaryCmd)
	glossaryCmd.AddCommand(glossaryListCmd)
	glossaryCmd.AddCommand(glossaryQueryCmd)
	glossaryCmd.AddCommand(glossaryAddCmd)

	glossaryListCmd.Flags().BoolVar(&glossaryJSON, "json", false, "Output glossary as JSON")
	glossaryQueryCmd.Flags().BoolVar(&glossaryJSON, "json", false, "Output matches as JSON")
	glossaryAddCmd.Flags().StringVar(&addDefinition, "definition", "", "Term definition (required)")
	glossaryAddCmd.Flags().StringSliceVar(&addAliases, "alias", nil, "Alias for the term (repeatable)")
	glossaryAddCmd.Flags().StringVar(&addCategory, "category", "", "Term category")
	glossaryAddCmd.Flags().StringSliceVar(&addSources, "source", nil, "Source marker (repeatable)")
}
