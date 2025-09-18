// Package helpers provides utility functions for OS detection and command execution.
package helpers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	log "github.com/charmbracelet/log"
	"github.com/charmbracelet/x/term"
)

var (
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	errorPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(lipgloss.Color("9")).
			Padding(0, 1).
			MarginLeft(2)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12"))
)

type outputMsg string

// RunCmd runs a bash command with a spinner and displays the output.
func RunCmd(command string) error {

	var cmd *exec.Cmd
	if strings.HasPrefix(command, "sudo") {
		pw, err := getSudoPassword()
		if err != nil {
			e := fmt.Errorf("failed to get sudo password: %w", err)
			PrintErrors(e)
			return e
		}
		cmd = exec.Command("sudo", "-S", "bash", "-c", command)
		cmd.Stdin = strings.NewReader(pw)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	go func() {
		for scanner.Scan() {
			fmt.Println(outputMsg(scanner.Text()))
		}
	}()

	err = cmd.Wait()

	return err
}

var (
	sudoPassword string
	sudoOnce     sync.Once
)

func getSudoPassword() (string, error) {
	var err error
	sudoOnce.Do(func() {
		fmt.Print("Enter sudo password: ")
		pw, e := term.ReadPassword(os.Stdin.Fd())
		fmt.Println()
		if e != nil {
			err = e
			return
		}
		sudoPassword = string(pw)
	})
	return sudoPassword, err
}

func PrintFail(msg string) {

	fmt.Println(errorStyle.Render("✗ " + msg))
}

func PrintErrors(err ...error) {

	PrintError("", err...)

}

func PrintError(msg string, err ...error) {
	PrintFail(msg)

	if len(err) > 0 && err[0] != nil {

		msg = fmt.Sprintf("%s: %v", msg, err[0])

	}
	log.Error(msg)
	fmt.Println(errorPanelStyle.Render("✗ " + msg))
}

func PrintSuccess(msg string) {
	log.Info(msg)

	fmt.Println(successStyle.Render("✓ " + msg))
}

func PrintInfo(msg string) {
	log.Debug(msg)
	fmt.Println(infoStyle.Render("ℹ " + msg))
}

// PromptYesNo prompts the user with a yes/no question and returns true for yes, false for no
func PromptYesNo(question string) bool {
	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to continue?").
				Value(&confirm),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Println("Error:", err)
		return false
	}

	if confirm {
		fmt.Println("Continuing...")
	} else {
		fmt.Println("Aborted.")
	}
	return confirm

}
