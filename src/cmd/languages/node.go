package languages

import (
	"fmt"
	"strings"

	"github.com/mpaxson/kettle/src/cmd/helpers"

	"github.com/spf13/cobra"
)

func installNVM() error {

	if helpers.CommandExists("nvm") {

		if !helpers.PromptYesNo("NVM is already installed. Do you want to reinstall it?") {
			return nil
		}
	}
	repo, err := helpers.GithubGetLatestRelease("nvm-sh", "nvm")
	if err != nil {
		helpers.PrintError("Failed to get latest NVM release", err)
		return err
	}
	tagVersion := *repo.TagName
	url := fmt.Sprintf("https://raw.githubusercontent.com/nvm-sh/nvm/%s/install.sh", tagVersion)
	helpers.PrintInfo(fmt.Sprintf("Downloading NVM install script from %q", url))
	if err := helpers.DownloadAndRunInstallScript(url, "nvm_install.sh"); err != nil {
		helpers.PrintError("Failed to install NVM", err)
		return err
	}

	return nil
}

func InstallNode() error {
	if !helpers.CommandExists("nvm") {
		err := installNVM()
		if err != nil {
			return err
		}
		helpers.PrintSuccess("NVM installed successfully.")

	}
	helpers.PrintInfo("Grabbing Latest Node.js...")
	// Add your npm update logic here
	err := helpers.RunCmdWithShellProfile("nvm install --lts")
	if err != nil {
		if strings.Contains(err.Error(), "is already installed") {
			helpers.PrintInfo("Latest node is already installed.")
		} else {
			helpers.PrintError("Failed to install Node.js via NVM", err)
			return err
		}

	}

	err = helpers.RunCmdWithShellProfile("nvm use --lts")
	if err != nil {
		helpers.PrintError("Failed to update Node.js via NVM", err)
		return err

	}
	helpers.PrintSuccess("Node.js updated successfully.")
	installNVMPath()
	return nil
}
func installNVMPath() {
	scriptPath := `
if command -v npm >/dev/null; then
    NPM_PREFIX=$(npm config get prefix 2>/dev/null)
    if [ -n "$NPM_PREFIX" ]; then
        export PATH="$NPM_PREFIX/bin:$PATH"
    fi
fi
`
	helpers.AddLineToKettleShellProfile(scriptPath)

}

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Node.js and npm related tools",
}

var nodeInstallCmd = &cobra.Command{
	Use:   "node",
	Short: "Install Node.js",
	Long:  `Downloads and installs Node.js.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := InstallNode()
		if err != nil {
			helpers.PrintError("Failed to install Node.js", err)
		} else {
			helpers.PrintSuccess("Node.js installed.")
		}
	},
}

var nvmInstallCmd = &cobra.Command{
	Use:   "nvm",
	Short: "Install NVM (Node Version Manager)",
	Long:  `Downloads and installs NVM, which allows you to install and manage multiple versions of Node.js.`,
	Run: func(cmd *cobra.Command, args []string) {
		helpers.PrintInfo("Installing NVM...")

		err := installNVM()
		if err != nil {
			helpers.PrintError("Failed to install NVM", err)
		} else {
			helpers.PrintSuccess("NVM installed.")
		}
	},
}

func init() {
	nodeCmd.AddCommand(nvmInstallCmd)
	nodeCmd.AddCommand(nodeInstallCmd)
}
