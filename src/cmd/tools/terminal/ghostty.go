package terminal

import (
	"fmt"
	"os/exec"

	"github.com/mpaxson/kettle/src/cmd/helpers"
	"github.com/spf13/cobra"
)

// GhosttyCmd represents the ghostty command
var GhosttyCmd = &cobra.Command{
	Use:   "ghostty",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var ghosttyInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs ghostty",
	Long:  `Installs ghostty for the current operating system.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !helpers.IsUbuntu() && !helpers.IsDarwin() {
			fmt.Println("Unsupported OS. Only Ubuntu and macOS are supported.")
			return
		}

		if _, err := exec.LookPath("ghostty"); err == nil {
			fmt.Println("ghostty is already installed.")
			return
		}

		if helpers.IsDarwin() {
			if _, err := exec.LookPath("brew"); err != nil {
				fmt.Println("Homebrew is not installed. Cannot install ghostty.")
				return
			}
			err := helpers.RunCmd("brew install ghostty")
			if err != nil {
				fmt.Println("Failed to install ghostty with brew:", err)
			}
		} else if helpers.IsUbuntu() {
			err := helpers.RunCmd("sudo snap install ghostty --classic")
			if err != nil {
				fmt.Println("Failed to install ghostty with snap:", err)
			}
		}
	},
}

var ghosttyBindF1Cmd = &cobra.Command{
	Use:   "bind-f1",
	Short: "Bind F1 to toggle Ghostty",
	Long:  `Bind F1 to toggle Ghostty.`,
	Run: func(cmd *cobra.Command, args []string) {
		const schema = "org.gnome.settings-daemon.plugins.media-keys"
		const path = "/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/"
		const script = "~/bin/ghostty-toggle.sh"
		const name = "Ghostty Toggle"
		const key = "F1"

		err := helpers.RunCmd(fmt.Sprintf("gsettings set %s custom-keybindings \"['%s']\"", schema, path))
		if err != nil {
			fmt.Println("Failed to set custom-keybindings:", err)
			return
		}

		err = helpers.RunCmd(fmt.Sprintf("gsettings set %s.custom-keybinding:%s name '%s'", schema, path, name))
		if err != nil {
			fmt.Println("Failed to set name:", err)
			return
		}

		err = helpers.RunCmd(fmt.Sprintf("gsettings set %s.custom-keybinding:%s command '%s'", schema, path, script))
		if err != nil {
			fmt.Println("Failed to set command:", err)
			return
		}

		err = helpers.RunCmd(fmt.Sprintf("gsettings set %s.custom-keybinding:%s binding '%s'", schema, path, key))
		if err != nil {
			fmt.Println("Failed to set binding:", err)
			return
		}

		fmt.Printf("✅ Bound %s to %s\n", key, script)
	},
}

var ghosttyUnbindF1Cmd = &cobra.Command{
	Use:   "unbind-f1",
	Short: "Remove F1 binding for Ghostty",
	Long:  `Remove F1 binding for Ghostty.`,
	Run: func(cmd *cobra.Command, args []string) {
		const schema = "org.gnome.settings-daemon.plugins.media-keys"
		const path = "/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/"
		const key = "F1"
		const script = "~/bin/ghostty-toggle.sh"

		err := helpers.RunCmd(fmt.Sprintf("gsettings reset %s custom-keybindings", schema))
		if err != nil {
			fmt.Println("Failed to reset custom-keybindings:", err)
			return
		}

		err = helpers.RunCmd(fmt.Sprintf("gsettings reset-recursively %s.custom-keybinding:%s", schema, path))
		if err != nil {
			fmt.Println("Failed to reset-recursively:", err)
			return
		}

		fmt.Printf("🗑️  Unbound %s from %s\n", key, script)
	},
}

func init() {
	GhosttyCmd.AddCommand(ghosttyInstallCmd)
	GhosttyCmd.AddCommand(ghosttyBindF1Cmd)
	GhosttyCmd.AddCommand(ghosttyUnbindF1Cmd)
}
