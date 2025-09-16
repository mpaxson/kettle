package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mpaxson/kettle/src/cmd/helpers"
	"github.com/spf13/cobra"
)

func generateCompletionForShell(shell string) error {
	configDir, err := helpers.GetKettleConfigDir()
	if err != nil {
		return err
	}
	completionsDir := filepath.Join(configDir, "completions")

	// Create the completions directory if it doesn't exist
	if err := os.MkdirAll(completionsDir, 0755); err != nil {
		return fmt.Errorf("could not create completions directory: %w", err)
	}

	// Create the completion file
	completionFile := filepath.Join(completionsDir, fmt.Sprintf("kettle.%s", shell))
	file, err := os.Create(completionFile)
	if err != nil {
		return fmt.Errorf("could not create completion file: %w", err)
	}
	defer file.Close()

	// Generate the completion script
	switch shell {
	case "bash":
		if err := rootCmd.GenBashCompletion(file); err != nil {
			return err
		}
	case "zsh":
		if err := rootCmd.GenZshCompletion(file); err != nil {
			return err
		}
	case "fish":
		if err := rootCmd.GenFishCompletion(file, true); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported shell for completion: %s", shell)
	}

	helpers.PrintSuccess(fmt.Sprintf("Generated %s completion file at %s", shell, completionFile))
	return nil
}

// GenerateAllCompletionFiles creates completion files for all supported shells.
func GenerateAllCompletionFiles() {
	shells := []string{"bash", "zsh", "fish"}
	for _, shell := range shells {
		if err := generateCompletionForShell(shell); err != nil {
			helpers.PrintError(fmt.Sprintf("Failed to generate %s completion", shell), err)
		}
	}
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install kettle to your system",
	Long:  `Install the kettle binary to a directory in your PATH.`,
	Run: func(cmd *cobra.Command, args []string) {
		exePath, err := os.Executable()
		if err != nil {
			return
		}
		err = helpers.InstallBinary(exePath)
		if err == nil {
			helpers.PrintSuccess("kettle installed successfully!")
		} else {
			helpers.PrintError("kettle installation failed:", err)
		}

		GenerateAllCompletionFiles()
		if err := helpers.EnsureCompletionsSourced(); err != nil {
			helpers.PrintError("Failed to source completion files", err)
		}

		err = helpers.EnsureKettleProfileSourced()
		if err != nil {
			helpers.PrintError("kettle profile sourcing failed:", err)
		} else {
			helpers.PrintSuccess("kettle profile sourced successfully!")
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
