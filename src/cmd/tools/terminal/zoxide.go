package terminal

import (
	"fmt"

	"github.com/mpaxson/kettle/src/cmd/helpers"
	"github.com/spf13/cobra"
)

// ZoxideCmd represents the zoxide command
var ZoxideCmd = &cobra.Command{
	Use:   "zoxide",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var zoxideInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs zoxide",
	Long:  `Installs zoxide.`,
	Run: func(cmd *cobra.Command, args []string) {
		if helpers.CommandExists("zoxide") {
			fmt.Println("Zoxide is already installed.")
			return
		}

		err := helpers.RunCmd("curl -sSfL https://raw.githubusercontent.com/ajeetdsouza/zoxide/main/install.sh | sh")
		if err != nil {
			helpers.PrintFail("Failed to install zoxide")
		} else {
			helpers.PrintSuccess("Zoxide installed.")
		}
		helpers.GetKettleShellProfile()
	},
}

func init() {
	ZoxideCmd.AddCommand(zoxideInstallCmd)
}
