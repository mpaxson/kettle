// Package tools provides commands for managing various development tools.
package tools

import (
	"github.com/mpaxson/kettle/src/cmd/tools/terminal"
	"github.com/spf13/cobra"
)

// ToolsCmd represents the tools command
var ToolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	ToolsCmd.AddCommand(terminal.TerminalCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ToolsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ToolsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
