# GitHub Copilot Instructions

## Goal

Migrate all Python Invoke tasks defined in `src/scripts/` into equivalent Go [Cobra](https://github.com/spf13/cobra) commands.

The existing Python tasks are structured as nested Invoke collections.  
Example:

```bash
inv terminal.ghostty.install
inv tools.terminal.ghostty.bind-f1
```

We want these to be turned into a Go CLI with the same hierarchical command layout using Cobra.

---

## Guidelines for Copilot

### General

- Always generate Go code using Cobra (`github.com/spf13/cobra`).
- Each Python Invoke namespace/collection (e.g., `terminal`, `tools`) should map to a Cobra command group.
- Each Invoke task (final leaf) should map to a Cobra subcommand.
- Commands must be defined under `cmd/` directory with `init()` to add them to the root or parent command.
- Use idiomatic Go CLI structure:

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

- Preserve the nesting of Invoke commands:

  - `inv terminal.ghostty.install` → `./mycli terminal ghostty install`
  - `inv tools.terminal.ghostty.bind-f1` → `./mycli tools terminal ghostty bind-f1`

- If an Invoke task uses `-` in the name, keep it as-is (`bind-f1`).

### Command Behavior

- For now, commands can `fmt.Println("running <command>")` as stubs.
- Later, they may execute corresponding shell logic (using `os/exec`), but initial focus is **scaffolding**.

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

- **Scaffold parent command**

  > "Create a Cobra command for `terminal` with description 'Terminal related commands'. Add it to `rootCmd`."

- **Add subcommand**

  > "Under `terminal`, add a `ghostty` command grouping Ghostty terminal commands."

- **Add task command**

  > "Add a subcommand `install` under `ghostty` that prints 'running terminal ghostty install'."

- **Repeat for all Invoke tasks in `src/scripts/`.**

---

## Notes

- Keep Go files small and modular (one command per file).
- Always attach new commands in `init()` so the CLI autowires.
- Follow naming conventions: lowercase, hyphenated for flags/commands if needed.
- Root command should describe the project purpose and show available subcommands.

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

### Helper Functions

For consistency, use the following helper functions from the `helpers` package when applicable.

- **`helpers.CommandExists("command-name")`**:

  - Use this to check if a command-line tool is available in the user's `PATH` before attempting to use it.
  - **Example**:
    ```go
    if !helpers.CommandExists("brew") {
        helpers.PrintError("Homebrew is not installed.")
        return
    }
    ```

- **`helpers.PrintError("Your error message")`**:

  - Use this for printing all user-facing error messages. It styles the output in red with a `✗` prefix.
  - **Example**:
    ```go
    if err != nil {
        helpers.PrintError("Failed to complete the task.")
    }
    ```

- **`helpers.PrintSuccess("Your success message")`**:

  - Use this for printing success messages. It styles the output in green with a `✓` prefix.
  - **Example**:
    ```go
    helpers.PrintSuccess("Task completed successfully.")
    ```

- **`helpers.PrintInfo("Your info message")`**:
  - Use this for printing informational messages. It styles the output in blue with a `ℹ` prefix.
  - **Example**:
    ```go
    helpers.PrintInfo("Downloading Go...")
    ```

### Shell Profile Management

For commands that need to modify shell configurations (like adding exports to PATH), use the following helper functions from the `helpers` package:

- **`helpers.GetCurrentShell()`**:

  - Returns the name of the currently running shell (bash, zsh, fish).
  - **Example**:
    ```go
    shell := helpers.GetCurrentShell()
    // Returns: "zsh", "bash", or "fish"
    ```

- **`helpers.GetShellProfile()`**:

  - Returns the path to the main shell configuration file based on the current shell.
  - **Example**:
    ```go
    shellProfile  := helpers.GetShellProfile()
    // Returns: "/home/user/.zshrc", "/home/user/.bashrc", etc.
    ```

- **`helpers.AddLineToShellProfile(line string)`**:

  - Adds a line to the user's main shell profile if it doesn't already exist.
  - **Example**:
    ```go
    helpers.AddLineToShellProfile("export PATH=$PATH:/usr/local/go/bin")
    ```

- **`helpers.GetKettleConfigDir()`**:

  - Returns the path to `~/.config/kettle` and creates it if it doesn't exist.
  - **Example**:
    ```go
    configDir, err := helpers.GetKettleConfigDir()
    // Returns: "/home/user/.config/kettle"
    ```

- **`helpers.GetKettleShellProfile()`**:

  - Returns the path to the kettle-specific shell profile file (e.g., `~/.config/kettle/kettle.zshrc`).
  - **Example**:
    ```go
    profilePath, err := helpers.GetKettleShellProfile()
    // Returns: "/home/user/.config/kettle/kettle.zshrc"
    ```

- **`helpers.AddLineToKettleShellProfile(line string)`**:

  - Adds a line to the kettle-specific shell profile if it doesn't already exist.
  - **Example**:
    ```go
    err := helpers.AddLineToKettleShellProfile("alias ll='ls -la'")
    ```

- **`helpers.EnsureKettleProfileSourced()`**:

  - Ensures the main shell profile sources the kettle-specific profile.
  - **Example**:
    ```go
    helpers.EnsureKettleProfileSourced()
    ```

- **`helpers.EnsureCompletionsSourced()`**:
  - Adds the source command for shell completions to the kettle profile.
  - **Example**:
    ```go
    err := helpers.EnsureCompletionsSourced()
    ```

### Best Practices for Shell Configuration

1. **Use Kettle-Specific Profiles**: Instead of directly modifying `~/.zshrc` or `~/.bashrc`, add configurations to the kettle-specific profile files in `~/.config/kettle/`.

2. **Check Before Adding**: Always use the helper functions that check if a line already exists before adding it to prevent duplicates.

3. **Inform Users**: After modifying shell profiles, always inform users that they need to restart their shell or source the profile for changes to take effect.

**Example Implementation**:

```go
// Add Go to PATH
if !helpers.AddLineToShellProfile("export PATH=$PATH:/usr/local/go/bin") {
    helpers.PrintError("Failed to update shell profile", err)
    return
}
helpers.PrintSuccess("Go added to PATH. Please restart your shell.")
```

### File Handling and Resource Management

When working with files in Go, always use proper resource management patterns to ensure files are closed and errors are handled correctly:

- **Use `defer helpers.FileClose(file, &err)`** for proper error handling when closing files:

  ```go
  func processFile(filename string) (err error) {
      file, err := os.Open(filename)
      if err != nil {
          return fmt.Errorf("failed to open file: %w", err)
      }
      defer helpers.FileClose(file, &err)

      // Process the file...
      return nil
  }
  ```

- **Always check errors when opening files**:

  ```go
  file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
      return fmt.Errorf("could not open file for writing: %w", err)
  }
  defer helpers.FileClose(file, &err)
  ```

- **Use named return values** when using `helpers.FileClose` to allow the helper to modify the return error if file closing fails:

  ```go
  func writeToFile(path, content string) (err error) {
      file, err := os.Create(path)
      if err != nil {
          return err
      }
      defer helpers.FileClose(file, &err)

      _, err = file.WriteString(content)
      return err
  }
  ```

**Why use `helpers.FileClose`?**

- Ensures file close errors are not silently ignored (violates errcheck linting)
- Provides consistent error handling across the codebase
- Automatically wraps close errors with context
- Prevents resource leaks even when errors occur

**Example Helper Implementation**:

```go
// FileClose safely closes a file and updates the error pointer if closing fails
func FileClose(file *os.File, err *error) {
    if closeErr := file.Close(); closeErr != nil {
        if *err == nil {
            *err = fmt.Errorf("failed to close file: %w", closeErr)
        } else {
            *err = fmt.Errorf("%w (also failed to close file: %v)", *err, closeErr)
        }
    }
}
```

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
```

## Example: Cobra Command

```go
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
```

## Copilot Prompts

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

