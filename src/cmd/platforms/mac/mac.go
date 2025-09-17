// Package mac contains macOS-specific implementations.
package mac

import (
	"github.com/spf13/cobra"
)

// MacCmd represents the macOS command
var MacCmd = &cobra.Command{
	Use:   "mac",
	Short: "Mac OS specific commands",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// LinuxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// LinuxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
