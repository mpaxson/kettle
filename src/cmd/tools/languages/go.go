package languages

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/mpaxson/kettle/src/cmd/helpers"
	"github.com/spf13/cobra"
)

func addGoToPath() {

	// Add ~/go/bin to PATH for Go binaries installed with 'go install'
	helpers.PrintInfo("Adding Go workspace bin to PATH...")
	if helpers.AddLineToKettleShellProfile(`export PATH=$PATH:$HOME/go/bin`) {
		helpers.PrintSuccess("Added Go workspace bin to kettle shell profile")
	}

	// Ensure kettle profile is sourced
	helpers.EnsureKettleProfileSourced()
	// Add to PATH
	helpers.PrintInfo("Adding Go to PATH...")
	if !helpers.AddLineToShellProfile(`export PATH=$PATH:/usr/local/go/bin`) {
		helpers.PrintInfo("Go already in PATH")
		return
	}
	helpers.PrintSuccess("Added Go to PATH")

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

	shellinfo := helpers.GetShellInfo()

	added := helpers.AddToProfileIfCmdExists(fmt.Sprintf(`eval "$(golangci-lint completion %s)"`, shellinfo.Type), "golangci-lint")
	if added {
		helpers.PrintSuccess("Added golangci-lint completions to shell profile")
	}
}
func installGoLint() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		helpers.PrintError("Failed to get user home directory", err)
		return
	}
	if gopath := os.Getenv("GOPATH"); gopath == "" {
		gopath = fmt.Sprintf("%s/go", homeDir)
		err := os.Setenv("GOPATH", gopath)
		if err != nil {
			helpers.PrintError("Failed to set GOPATH", err)
			return
		}
	}

	helpers.PrintInfo("Downloading golangci-lint install script...")
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = filepath.Join(homeDir, "go")
	}

	installDir := filepath.Join(goPath, "bin")
	helpers.AddToPath(installDir)
	destPath, err := helpers.GithubDownloadLatestRelease("golangci", "golangci-lint", installDir, "golangci-lint")
	helpers.PrintInfo("Downloaded golangci-lint to: " + destPath)

	if err != nil {
		log.Error(destPath)
		helpers.PrintError("Failed to install golangci-lint", err)
		return
	}

	helpers.PrintSuccess("golangci-lint installed successfully.")
	addGoLintToPath()

}

var goLintInstallCmd = &cobra.Command{
	Use:   "golangci-lint",
	Short: "Install lint tool for Go",
	Long:  `Downloads and installs the latest version of golangci-lint.`,
	Run: func(cmd *cobra.Command, args []string) {
		if helpers.CommandExists("golangci-lint") {
			// Prompt user if they want to reinstall
			if !helpers.PromptYesNo("golangci-lint is already installed. Do you want to reinstall it?") {
				helpers.PrintInfo("Skipping golangci-lint installation.")
				addGoLintToPath()
				return
			}
			helpers.PrintInfo("Proceeding with golangci-lint reinstallation...")
		}
		installGoLint()
	},
}

func init() {
	goCmd.AddCommand(goInstallCmd)
	goCmd.AddCommand(goLintInstallCmd)

	// You'll need to add goCmd to the parent 'languages' command.
	// Example: LanguagesCmd.AddCommand(goCmd)
}
