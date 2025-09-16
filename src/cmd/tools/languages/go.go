package languages

import (
	"fmt"
	"runtime"

	"github.com/mpaxson/kettle/src/cmd/helpers"
	"github.com/spf13/cobra"
)

func addGoToPath() {

	// Add ~/go/bin to PATH for Go binaries installed with 'go install'
	helpers.PrintInfo("Adding Go workspace bin to PATH...")
	if err := helpers.AddLineToKettleShellProfile(`export PATH=$PATH:$HOME/go/bin`); err != nil {
		helpers.PrintError("Failed to update kettle shell profile", err)
		return
	}

	// Ensure kettle profile is sourced
	if err := helpers.EnsureKettleProfileSourced(); err != nil {
		helpers.PrintError("Failed to source kettle profile", err)
		return
	}

	// Add to PATH
	helpers.PrintInfo("Adding Go to PATH...")
	if err := helpers.AddLineToShellProfile(`export PATH=$PATH:/usr/local/go/bin`); err != nil {
		helpers.PrintError("Failed to update shell profile", err)
		return
	}

}

var goCmd = &cobra.Command{
	Use:   "go",
	Short: "Install Go programming language",
}

var goInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Go",
	Long:  `Downloads and installs the latest version of Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		if helpers.CommandExists("go") {
			addGoToPath()
			return

		}

		version := "latest" // It's good practice to manage this version string
		arch := runtime.GOARCH
		goTarball := fmt.Sprintf("go%s.linux-%s.tar.gz", version, arch)
		downloadURL := fmt.Sprintf("https://golang.org/dl/%s", goTarball)

		// Download
		helpers.PrintInfo(fmt.Sprintf("Downloading Go %s...", version))
		if err := helpers.RunCmd(fmt.Sprintf("curl -OL %s", downloadURL)); err != nil {
			helpers.PrintError("Failed to download Go tarball", err)
			return
		}

		// Install
		helpers.PrintInfo("Installing Go...")
		installCommands := []string{
			"sudo rm -rf /usr/local/go",
			fmt.Sprintf("sudo tar -C /usr/local -xzf %s", goTarball),
			fmt.Sprintf("rm %s", goTarball),
		}

		for _, command := range installCommands {
			if err := helpers.RunCmd(command); err != nil {
				helpers.PrintError(fmt.Sprintf("Failed to execute command: %s", command), err)
				return
			}
		}
		addGoToPath()

	},
}

func addGoLintToPath() {

	// Add ~/go/bin to PATH for Go binaries installed with 'go install'
	helpers.PrintInfo("Adding Go langci completions to path")

	shellinfo, err := helpers.GetShellInfo()
	if err != nil {
		helpers.PrintError("Failed to get shell info", err)
		return
	}
	if err := helpers.AddToProfileIfCmdExists(fmt.Sprintf(`eval "$(golangci-lint completion %s)"`, shellinfo.Type), "golangci-lint"); err != nil {
		helpers.PrintError("Failed to update kettle shell profile", err)
		return
	}
	helpers.PrintSuccess("Added golangci-lint completions to shell profile")
}

var goLintInstallCmd = &cobra.Command{
	Use:   "golangci-lint",
	Short: "Install lint tool for Go",
	Long:  `Downloads and installs the latest version of golangci-lint.`,
	Run: func(cmd *cobra.Command, args []string) {
		if helpers.CommandExists("golangci-lint") {
			addGoLintToPath()
			return

		}
		err := helpers.RunCmd("curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.4.0")
		if err != nil {
			helpers.PrintError("Failed to install golangci-lint", err)
			return
		}
		helpers.PrintSuccess("golangci-lint installed successfully.")
		addGoLintToPath()

	},
}

func init() {
	goCmd.AddCommand(goInstallCmd)
	goCmd.AddCommand(goLintInstallCmd)

	// You'll need to add goCmd to the parent 'languages' command.
	// Example: LanguagesCmd.AddCommand(goCmd)
}
