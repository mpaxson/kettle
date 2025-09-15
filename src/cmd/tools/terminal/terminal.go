package terminal

import (
	"github.com/spf13/cobra"
)

// TerminalCmd represents the terminal command
var TerminalCmd = &cobra.Command{
	Use:   "terminal",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	TerminalCmd.AddCommand(GhosttyCmd)
	TerminalCmd.AddCommand(KittyCmd)
	TerminalCmd.AddCommand(AutoenvCmd)
	TerminalCmd.AddCommand(ZoxideCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// TerminalCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// TerminalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
