package terminal

import (
	"fmt"

	"github.com/mpaxson/kettle/src/cmd/helpers"
	"github.com/spf13/cobra"
)

// KittyCmd represents the kitty command
var KittyCmd = &cobra.Command{
	Use:   "kitty",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var kittyInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs Kitty",
	Long:  `Installs Kitty.`,
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{
			"curl -L https://sw.kovidgoyal.net/kitty/installer.sh | sh /dev/stdin",
			"ln -sf ~/.local/kitty.app/bin/kitty ~/.local/kitty.app/bin/kitten ~/.local/bin/",
			"cp ~/.local/kitty.app/share/applications/kitty.desktop ~/.local/share/applications/",
			"cp ~/.local/kitty.app/share/applications/kitty-open.desktop ~/.local/share/applications/",
			`sed -i "s|Icon=kitty|Icon=/home/$USER/.local/kitty.app/share/icons/hicolor/256x256/apps/kitty.png|g" ~/.local/share/applications/kitty*.desktop`,
			`sed -i "s|Exec=kitty|Exec=/home/$USER/.local/kitty.app/bin/kitty|g" ~/.local/share/applications/kitty*.desktop`,
		}

		for _, command := range commands {
			err := helpers.RunCmd(command)
			if err != nil {
				fmt.Printf("Failed to execute command: %s\n", command)
				fmt.Println(err)
			}
		}
	},
}

func init() {
	KittyCmd.AddCommand(kittyInstallCmd)
}
