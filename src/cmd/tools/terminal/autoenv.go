package terminal

// autoenv.go
import (
	"fmt"
	"os/exec"

	"github.com/mpaxson/kettle/src/cmd/helpers"
	"github.com/spf13/cobra"
)

// AutoenvCmd represents the autoenv command
var AutoenvCmd = &cobra.Command{
	Use:   "autoenv",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var autoenvInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs autoenv",
	Long:  `Installs autoenv.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := exec.LookPath("autoenv"); err == nil {
			fmt.Println("autoenv is already installed.")
			return
		}

		if _, err := exec.LookPath("npm"); err == nil {
			err := helpers.RunCmd("npm install -g autoenv")
			if err != nil {
				fmt.Println("Failed to install autoenv with npm:", err)
			}
			return
		}

		if helpers.IsDarwin() {
			if _, err := exec.LookPath("brew"); err != nil {
				fmt.Println("Homebrew is not installed. Cannot install autoenv.")
				return
			}
			err := helpers.RunCmd("brew install autoenv")
			if err != nil {
				fmt.Println("Failed to install autoenv with brew:", err)
			}
		} else if helpers.IsUbuntu() {
			err := helpers.RunCmd("sudo apt install -y autoenv")
			if err != nil {
				fmt.Println("Failed to install autoenv with apt:", err)
			}
		} else {
			fmt.Println("Unsupported OS. Only Ubuntu and macOS are supported.")
		}
	},
}

func init() {
	AutoenvCmd.AddCommand(autoenvInstallCmd)
}
