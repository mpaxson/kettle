// Package languages provides commands for installing and managing programming languages.
package languages

import (
	"github.com/spf13/cobra"
)

// LanguagesCmd represents the languages command group.
var LanguagesCmd = &cobra.Command{
	Use:   "languages",
	Short: "Commands for installing and managing programming languages",
}

func init() {
	// Register the 'go' command as a subcommand of 'languages'
	LanguagesCmd.AddCommand(goCmd)
}
