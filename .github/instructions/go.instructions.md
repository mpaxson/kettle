
# GitHub Copilot Instructions

## Goal
Migrate all Python Invoke tasks defined in `src/scripts/` into equivalent Go [Cobra](https://github.com/spf13/cobra) commands.

The existing Python tasks are structured as nested Invoke collections.  
Example:
```bash
inv terminal.ghostty.install
inv tools.terminal.ghostty.bind-f1
````

We want these to be turned into a Go CLI with the same hierarchical command layout using Cobra.

---

## Guidelines for Copilot

### General

* Always generate Go code using Cobra (`github.com/spf13/cobra`).
* Each Python Invoke namespace/collection (e.g., `terminal`, `tools`) should map to a Cobra command group.
* Each Invoke task (final leaf) should map to a Cobra subcommand.
* Commands must be defined under `cmd/` directory with `init()` to add them to the root or parent command.
* Use idiomatic Go CLI structure:

  ```
  /cmd
    root.go
    terminal.go
    ghostty.go
    install.go
    bindf1.go
  /main.go
  ```

### Command Naming

* Preserve the nesting of Invoke commands:

  * `inv terminal.ghostty.install` → `./mycli terminal ghostty install`
  * `inv tools.terminal.ghostty.bind-f1` → `./mycli tools terminal ghostty bind-f1`
* If an Invoke task uses `-` in the name, keep it as-is (`bind-f1`).

### Command Behavior

* For now, commands can `fmt.Println("running <command>")` as stubs.
* Later, they may execute corresponding shell logic (using `os/exec`), but initial focus is **scaffolding**.

### Examples

#### Python Invoke

```python
# src/scripts/terminal/ghostty.py
@task
def install(c):
    """Install Ghostty terminal"""
    ...

@task
def bind_f1(c):
    """Bind F1 to toggle Ghostty terminal"""
    ...
```

#### Go Cobra

```go
// cmd/terminal.go
var terminalCmd = &cobra.Command{
    Use:   "terminal",
    Short: "Terminal related commands",
}

func init() {
    rootCmd.AddCommand(terminalCmd)
}

// cmd/ghostty.go
var ghosttyCmd = &cobra.Command{
    Use:   "ghostty",
    Short: "Ghostty terminal commands",
}

func init() {
    terminalCmd.AddCommand(ghosttyCmd)
}

// cmd/install.go
var installCmd = &cobra.Command{
    Use:   "install",
    Short: "Install Ghostty terminal",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("running terminal ghostty install")
    },
}

func init() {
    ghosttyCmd.AddCommand(installCmd)
}

// cmd/bindf1.go
var bindF1Cmd = &cobra.Command{
    Use:   "bind-f1",
    Short: "Bind F1 to toggle Ghostty",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("running tools terminal ghostty bind-f1")
    },
}

func init() {
    ghosttyCmd.AddCommand(bindF1Cmd)
}
```

---

## Copilot Prompts

When editing or creating files under `cmd/`, follow this workflow:

* **Scaffold parent command**

  > "Create a Cobra command for `terminal` with description 'Terminal related commands'. Add it to `rootCmd`."

* **Add subcommand**

  > "Under `terminal`, add a `ghostty` command grouping Ghostty terminal commands."

* **Add task command**

  > "Add a subcommand `install` under `ghostty` that prints 'running terminal ghostty install'."

* **Repeat for all Invoke tasks in `src/scripts/`.**

---

## Notes

* Keep Go files small and modular (one command per file).
* Always attach new commands in `init()` so the CLI autowires.
* Follow naming conventions: lowercase, hyphenated for flags/commands if needed.
* Root command should describe the project purpose and show available subcommands.



---

# GitHub Copilot Instructions

## Goal
Use [Charmbracelet](https://charm.sh/) libraries (`bubbletea`, `bubbles`, `lipgloss`) to make terminal UX more polished when running Cobra commands.  

Every command should:
- Show a **spinner** (from `bubbles/spinner`) while long tasks (like `inv` or `exec.Command`) are running.
- Print **Lipgloss-styled success/error/info messages** after completion.

---

## Libraries to Use
- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/bubbles/spinner`
- `github.com/charmbracelet/lipgloss`

---

## Guidelines for Copilot (Terminal Output )

### Spinner
- Always initialize a spinner with `spinner.New()` and set its style.
- Use `spinner.Dot` or `spinner.Line` as defaults.
- Update the spinner on `tea.TickMsg`.
- Show `"Running task..."` while a command is executing.

### Styled Messages
- Use `lipgloss.NewStyle()` for success/error/info text.
- Color conventions:
  - ✅ Success → Green (`Foreground(lipgloss.Color("10"))`)
  - ❌ Error → Red (`Foreground(lipgloss.Color("9"))`)
  - ℹ️ Info / spinner → Blue (`Foreground(lipgloss.Color("12"))`)

### Command Integration
- Commands should return a `RunE` function in Cobra that calls a helper like `RunWithSpinner("task description", func() error { ... })`.

---

## Example: Spinner + Success/Error

```go
package ui

import (
	"fmt"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
)

type model struct {
	spinner spinner.Model
	msg     string
	done    bool
	err     error
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.TickMsg:
		if m.done {
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.done {
		if m.err != nil {
			return errorStyle.Render("✗ " + m.err.Error())
		}
		return successStyle.Render("✓ " + m.msg)
	}
	return infoStyle.Render(m.spinner.View() + " " + m.msg)
}

// RunWithSpinner runs a shell command with spinner feedback
func RunWithSpinner(msg string, cmd *exec.Cmd) error {
	s := spinner.New()
	s.Spinner = spinner.Dot
	m := model{spinner: s, msg: msg}

	p := tea.NewProgram(m)
	go func() {
		err := cmd.Run()
		m.done = true
		m.err = err
		if err == nil {
			m.msg = msg + " completed"
		}
		p.Send(tea.QuitMsg{}) // quit after task finishes
	}()
	_, err := p.Run()
	return err
}
Example: Cobra Command
go
Copy code
var ghosttyInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Ghostty terminal",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ui.RunWithSpinner(
			"Installing Ghostty",
			exec.Command("inv", "terminal.ghostty.install"),
		)
	},
}
Copilot Prompts
When writing new commands:

Scaffold spinner helper

"Generate a Bubbletea spinner model with Lipgloss styles for success, error, and info."

Wrap an exec.Command

"Wrap exec.Command("inv", "tools.terminal.ghostty.bind-f1") in a Bubbletea spinner that shows 'Binding F1...' and prints green success or red error messages."

Add to Cobra

"Create a Cobra command bind-f1 under ghostty that calls ui.RunWithSpinner to run inv tools.terminal.ghostty.bind-f1."

Notes
Always use RunE in Cobra commands to bubble up errors.

Spinner should quit automatically when the task ends.

Use Lipgloss styles consistently across all commands for a professional look.

yaml
Copy code
