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

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type model struct {
	spinner spinner.Model
	output  string
	done    bool
	err     error
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg: // ✅ use spinner.TickMsg, not tea.TickMsg
		if m.done {
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case outputMsg:
		m.output = string(msg)
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	if m.done {
		if m.err != nil {
			return errorStyle.Render("✗ " + m.err.Error())
		}
		return successStyle.Render("✓ Done!")
	}
	s := infoStyle.Render(m.spinner.View() + " Running task...")
	if m.output != "" {
		s += "\n" + m.output
	}
	return s
}

// RunCmd runs a bash command with a spinner and displays the output.
func RunCmd(command string) error {
	s := spinner.New()
	s.Spinner = spinner.Dot
	m := model{spinner: s}

	// Start the spinner program
	p := tea.NewProgram(m)
	go func() {
		var cmd *exec.Cmd
		if strings.HasPrefix(command, "sudo") {
			pw, err := getSudoPassword()
			if err != nil {
				m.err = err
				p.Send(tea.QuitMsg{})
				return
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
			m.err = err
			p.Send(tea.QuitMsg{})
			return
		}

		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		go func() {
			for scanner.Scan() {
				p.Send(outputMsg(scanner.Text()))
			}
		}()

		err = cmd.Wait()

		m.done = true
		m.err = err

		// force quit bubbletea loop
		p.Send(tea.QuitMsg{})
	}()

	if _, err := p.Run(); err != nil {
		return err
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
	fmt.Println(errorPanelStyle.Render("✗ " + msg))
}

func PrintSuccess(msg string) {
	fmt.Println(successStyle.Render("✓ " + msg))
}

func PrintInfo(msg string) {
	fmt.Println(infoStyle.Render("ℹ " + msg))
}
