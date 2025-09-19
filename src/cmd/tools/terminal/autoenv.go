package terminal

// autoenv.go
import (
	"strings"

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

func AddAutoenvToProfile() {
	script := `if command -v npm >/dev/null; then
  NPM_ROOT=$(npm root -g 2>/dev/null)
  if [ -n "$NPM_ROOT" ] && [ -f "$NPM_ROOT/@hyperupcall/autoenv/activate.sh" ]; then
    source "$NPM_ROOT/@hyperupcall/autoenv/activate.sh"
  fi
fi`

	helpers.AddLineToKettleShellProfile(script)
}

var autoenvInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs autoenv",
	Long:  `Installs autoenv.`,
	Run: func(cmd *cobra.Command, args []string) {

		if !helpers.CommandExists("git") {
			helpers.PrintFail("git is not installed. Cannot install autoenv.")
		}

		err := helpers.RunCmd("git clone 'https://github.com/hyperupcall/autoenv' ~/.autoenv", false)

		if err != nil && strings.Contains(err.Error(), "128") {
			helpers.PrintInfo("autoenv is already cloned.")
			err := helpers.RunCmd("cd ~/.autoenv; git pull", false)
			if err != nil {
				helpers.PrintError("Failed to update autoenv", err)
			} else {
				helpers.PrintSuccess("Updated autoenv.")
			}

		} else {
			helpers.PrintError("Failed to clone autoenv repository", err)
			panic(err)
		}

		helpers.AddLineToKettleShellProfile("source ~/.autoenv/activate.sh")
	},
}

func init() {
	AutoenvCmd.AddCommand(autoenvInstallCmd)
}
