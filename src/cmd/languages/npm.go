package languages

import (
	"fmt"

	"github.com/mpaxson/kettle/src/cmd/helpers"

	"github.com/spf13/cobra"
)

func installNvm() {

	if helpers.CommandExists("nvm") {

		if !helpers.PromptYesNo("NVM is already installed. Do you want to reinstall it?") {
			return
		}
	}
	repo, err := helpers.GithubGetLatestRelease("nvm-sh", "nvm")
	if err != nil {
		helpers.PrintError("Failed to get latest NVM release", err)
		return
	}
	tagVersion := *repo.TagName
	url := fmt.Sprintf("https://raw.githubusercontent.com/nvm-sh/nvm/%s/install.sh", tagVersion)
	helpers.PrintInfo(fmt.Sprintf("Downloading NVM install script from %q", url))
	if err := helpers.DownloadAndRunInstallScript(url, "nvm_install.sh"); err != nil {
		helpers.PrintError("Failed to install NVM", err)
		return
	}
	helpers.PrintSuccess("NVM installed successfully.")

}

var npmCmd = &cobra.Command{
	Use:   "npm",
	Short: "Node.js and npm related tools",
}

var nvmInstallCmd = &cobra.Command{
	Use:   "nvm",
	Short: "Install NVM (Node Version Manager)",
	Long:  `Downloads and installs NVM, which allows you to install and manage multiple versions of Node.js.`,
	Run: func(cmd *cobra.Command, args []string) {
		installNvm()

	},
}

func init() {
	npmCmd.AddCommand(nvmInstallCmd)
}
