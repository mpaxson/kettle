package terminal

// ghostty.go
import (
	"fmt"
	"os"
	"path/filepath"

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
			helpers.PrintFail("Unsupported OS. Only Ubuntu and macOS are supported.")
			return
		}

		if helpers.CommandExists("ghostty") {
			fmt.Println("ghostty is already installed.")
			return
		}

		if helpers.IsDarwin() {
			if !helpers.CommandExists("brew") {
				helpers.PrintFail("Homebrew is not installed. Cannot install ghostty.")
				return
			}
			err := helpers.RunCmd("brew install ghostty")
			if err != nil {
				helpers.PrintFail("Failed to install ghostty with brew")
			}
		} else if helpers.IsUbuntu() {
			err := helpers.RunCmd("sudo snap install ghostty --classic")
			if err != nil {
				helpers.PrintFail("Failed to install ghostty with snap")
			}
		}
	},
}
var ghosttyCreateToggleScriptCmd = &cobra.Command{
	Use:   "create-toggle-script",
	Short: "Creates the ghostty toggle script",
	Long:  `Creates the ghostty toggle script at ~/bin/ghostty-toggle.sh.`,
	Run: func(cmd *cobra.Command, args []string) {
		scriptPath := getScriptPath()
		const scriptContent = `#!/usr/bin/env bash
# Toggle Ghostty as a dropdown terminal (fullscreen on F1)
APP="ghostty"
WM_CLASS="Ghostty"
BASHCMD="bash -c tmux attach || tmux new"
CMD="ghostty --class=$WM_CLASS --fullscreen -e $BASHCMD"
# tdrop options:
# -ma   : match by WM_CLASS
# -w    : window width in %
# -h    : window height in %
# -y    : drop from top (y=0)
tdrop -ma -w 5000 -h 5000 -y 0 -a $CMD`

		// Create bin directory if it doesn't exist
		binDir := filepath.Dir(scriptPath)
		if err := os.MkdirAll(binDir, 0755); err != nil {
			helpers.PrintFail(fmt.Sprintf("Failed to create bin directory: %v", err))
			return
		}

		err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
		if err != nil {
			helpers.PrintFail(fmt.Sprintf("Failed to create toggle script: %v", err))
			return
		}

		helpers.PrintSuccess(fmt.Sprintf("✅ Created and chmod +x %s", scriptPath))
	},
}

func getScriptPath() string {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		helpers.PrintFail("Failed to get user home directory")
		return ""
	}

	scriptPath := filepath.Join(homeDir, "bin", "ghostty-toggle.sh")
	return scriptPath
}

var ghosttyBindF1Cmd = &cobra.Command{
	Use:   "bind-f1",
	Short: "Bind F1 to toggle Ghostty",
	Long:  `Bind F1 to toggle Ghostty.`,
	Run: func(cmd *cobra.Command, args []string) {
		schema := "org.gnome.settings-daemon.plugins.media-keys"
		path := "/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/"
		name := "Ghostty Toggle"
		key := "F1"
		scriptPath := getScriptPath()

		// Check if script exists
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			helpers.PrintFail(fmt.Sprintf("❌ Error: %s not found. Make sure toggle-ghostty exists.", scriptPath))
			return
		}

		// Check if script is executable
		if info, err := os.Stat(scriptPath); err == nil {
			if info.Mode()&0111 == 0 {
				helpers.PrintFail(fmt.Sprintf("❌ Error: %s is not executable. Run: chmod +x %s", scriptPath, scriptPath))
				return
			}
		}

		// Register the custom binding path - this sets the array of custom keybindings
		err := helpers.RunCmd(fmt.Sprintf("gsettings set %s custom-keybindings \"['%s']\"", schema, path))
		if err != nil {
			helpers.PrintFail("Failed to set custom-keybindings")
			return
		}

		// Configure the binding name
		err = helpers.RunCmd(fmt.Sprintf("gsettings set %s.custom-keybinding:%s name \"%s\"", schema, path, name))
		if err != nil {
			helpers.PrintFail("Failed to set binding name")
			return
		}

		// Configure the binding command
		err = helpers.RunCmd(fmt.Sprintf("gsettings set %s.custom-keybinding:%s command \"%s\"", schema, path, scriptPath))
		if err != nil {
			helpers.PrintFail("Failed to set binding command")
			return
		}

		// Configure the key binding
		err = helpers.RunCmd(fmt.Sprintf("gsettings set %s.custom-keybinding:%s binding \"%s\"", schema, path, key))
		if err != nil {
			helpers.PrintFail("Failed to set key binding")
			return
		}

		helpers.PrintSuccess(fmt.Sprintf("✅ Bound %s to %s", key, scriptPath))
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
		scriptPath := getScriptPath()

		err := helpers.RunCmd(fmt.Sprintf("gsettings reset %s custom-keybindings", schema))
		if err != nil {
			helpers.PrintFail("Failed to reset custom-keybindings")
			return
		}

		err = helpers.RunCmd(fmt.Sprintf("gsettings reset-recursively %s.custom-keybinding:%s", schema, path))
		if err != nil {
			helpers.PrintFail("Failed to reset-recursively")
			return
		}

		helpers.PrintSuccess(fmt.Sprintf("🗑️  Unbound %s from %s", key, scriptPath))
	},
}

func init() {
	GhosttyCmd.AddCommand(ghosttyInstallCmd)
	GhosttyCmd.AddCommand(ghosttyBindF1Cmd)
	GhosttyCmd.AddCommand(ghosttyUnbindF1Cmd)
	GhosttyCmd.AddCommand(ghosttyCreateToggleScriptCmd)

}
