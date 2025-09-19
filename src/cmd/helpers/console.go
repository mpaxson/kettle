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
	cmdStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#555555")). // dark grey background
			Foreground(lipgloss.Color("#AAAAAA")). // light text
			Padding(0, 2)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder(), true, true, true, true).
			BorderForeground(lipgloss.Color("10"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder(), true, true, true, true).
			BorderForeground(lipgloss.Color("9"))
	errorPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(lipgloss.Color("9")).
			Padding(0, 1).
			MarginLeft(2)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12"))
)


// RunCmd runs a bash command with TUI output and displays the results in panels.
// Optional output parameter controls whether to print command output (default: true)
func RunCmd(command string, output ...bool) error {
	// Determine if we should print output (default: true)
	shouldPrintOutput := true
	if len(output) > 0 {
		shouldPrintOutput = output[0]
	}

	// Start TUI if not already running
	// Update status panel to show current command
	if shouldPrintOutput {
		PrintInfo(fmt.Sprintf("Running: %s", command))
	}

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
		v := fmt.Errorf("failed to start command: %w", err)
		if shouldPrintOutput {
			PrintErrors(v)
		}
		return v
	}

	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	go func() {
		for scanner.Scan() {
			if shouldPrintOutput {
				PrintCmdOutput(scanner.Text())
			}
		}
	}()

	err = cmd.Wait()

	if err != nil {
		if shouldPrintOutput {
			PrintError(fmt.Sprintf("Command failed: %s", command))
		}
		return err

	}

	if shouldPrintOutput {
		PrintInfo(fmt.Sprintf("Completed: %s", command))
	}
	return nil
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
	fmt.Println(errorPanelStyle.Render("✗  " + msg))
}

func PrintSuccess(msg string) {
	log.Info(msg)

	fmt.Println(successStyle.Render("✓  " + msg))
}

func PrintInfo(msg string) {
	log.Debug(msg)
	fmt.Println(infoStyle.Render("ℹ  " + msg))
}

func PrintCmdOutput(msg string) {
	log.Debug(msg)
	fmt.Println(cmdStyle.Render(" ➜  CmdOut: " + msg))
}

// PromptYesNo prompts the user with a yes/no question and returns true for yes, false for no
func PromptYesNo(question string) bool {
	var confirm bool

	fmt.Println(question)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(question).
				Description("Do you want to continue?").
				Value(&confirm),
		),
	)

	if err := form.Run(); err != nil {
		PrintError("Error:", err)
		return false
	}

	if confirm {
		fmt.Println("Continuing...")
	} else {
		fmt.Println("Aborted.")
	}
	return confirm

}
