package cmd

import (
	"fmt"

	"github.com/mpaxson/kettle/src/cmd/helpers"
	"github.com/mpaxson/kettle/src/internal/version"
	"github.com/spf13/cobra"
)

// Version information set at build time
var (
	Version = version.Version   // e.g. "1.0.0"
	Commit  = version.Commit    // e.g. "abc1234"
	Date    = version.BuildDate // e.g. "2024-06-01T12:00:00Z"
)

// SetVersionInfo sets the version information from main package
func SetVersionInfo(version, commit, date string) {
	if version != "" {
		Version = version
	}
	if commit != "" {
		Commit = commit
	}
	if date != "" {
		Date = date
	}
}

// GetVersion returns the formatted version string
func GetVersion() string {
	if Version == "dev" {
		return fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, Date)
	}
	return fmt.Sprintf("v%s (commit: %s, built: %s)", Version, Commit, Date)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of kettle",
	Long:  `Display the current version, commit hash, and build date of kettle.`,
	Run: func(cmd *cobra.Command, args []string) {
		helpers.PrintInfo(fmt.Sprintf("kettle %s", GetVersion()))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Add version flag to root command (Cobra built-in pattern)
	rootCmd.Version = GetVersion()
}
