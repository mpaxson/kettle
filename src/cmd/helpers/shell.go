package helpers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ShellInfo struct {
	Type         string
	ShellBinPath string
	Path         string
	KettlePath   string
	KettleConfig string
}

var (
	cached OSRelease
	once   sync.Once
	err    error

	cachedShell ShellInfo
	shellOnce   sync.Once
	shellErr    error
)

// GetShellInfo returns cached shell information, determining it once
func GetShellInfo() (ShellInfo, error) {
	shellOnce.Do(func() {

		shellPath := os.Getenv("SHELL")
		if shellPath == "" {
			shellErr = fmt.Errorf("SHELL environment variable not set")
			return
		}

		shellType := filepath.Base(shellPath)
		homeDir, err := os.UserHomeDir()
		if err != nil {
			shellErr = fmt.Errorf("could not get user home directory: %w", err)
			return
		}
		var shellProfilePath string

		switch shellType {
		case "bash":
			shellProfilePath = filepath.Join(homeDir, ".bashrc")
		case "zsh":
			shellProfilePath = filepath.Join(homeDir, ".zshrc")
		case "fish":
			shellProfilePath = filepath.Join(homeDir, ".config", "fish", "config.fish")
		default:
			shellErr = fmt.Errorf("unsupported shell type: %s", shellType)
			return

		}
		// Get kettle config directory

		configDir, err := GetKettleConfigDir()
		if err != nil {
			shellErr = fmt.Errorf("could not get kettle config directory: %w", err)
			return
		}
		kettlePath := filepath.Join(configDir, fmt.Sprintf("kettle.%src", shellType))

		cachedShell = ShellInfo{
			Type:         shellType,
			ShellBinPath: shellPath,
			Path:         shellProfilePath,
			KettleConfig: configDir,
			KettlePath:   kettlePath,
		}
	})
	return cachedShell, shellErr
}

// GetCurrentShell determines the name of the currently running shell by inspecting the SHELL environment variable.
func GetCurrentShell() string {
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		return "" // Default or error case
	}
	return filepath.Base(shellPath)
}

// GetShellProfile returns the path to the shell configuration file based on the current shell.
// It supports bash, zsh, and fish.

// AddLineToShellProfile adds a given line of text to the appropriate shell profile file
// if it does not already exist in the file.
func AddLineToShellProfile(line string) error {
	shellInfo, err := GetShellInfo()
	if err != nil {
		return fmt.Errorf("could not get shell info: %w", err)
	}

	// Check if the line already exists in the file.
	exists, err := ExistsInFile(shellInfo.Path, line)
	if err != nil {
		return fmt.Errorf("error checking shell profile: %w", err)
	}
	if exists {
		PrintInfo(fmt.Sprintf("Configuration already exists in %s.", filepath.Base(shellInfo.Path)))
		return nil
	}

	// If the file doesn't exist or the line isn't in it, append the line.
	file, err := os.OpenFile(shellInfo.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open shell profile for writing: %w", err)
	}
	defer fileClose(file, &err)
	//defer file.Close()

	if _, err := fmt.Fprintln(file, line); err != nil {
		return fmt.Errorf("failed to write to shell profile: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Configuration added to %s.", filepath.Base(shellInfo.Path)))
	return nil
}

// ExistsInFile checks if a given string `content` exists within the file at `filePath`.
// It reads the file line by line to avoid loading large files into memory.
func ExistsInFile(filePath, content string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer fileClose(file, &err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), content) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func GetKettleConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}
	configDir := filepath.Join(homeDir, ".config", "kettle")

	// Create the directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return "", fmt.Errorf("could not create kettle config directory: %w", err)
		}
	}
	return configDir, nil
}

// EnsureKettleProfileSourced makes sure the main shell profile sources the kettle-specific profile.
func EnsureKettleProfileSourced() error {
	shellInfo, err := GetShellInfo()
	if err != nil {
		return err
	}

	sourceCmd := fmt.Sprintf("source %s", shellInfo.KettlePath)

	return AddLineToShellProfile(sourceCmd)
}
func addToFile(input string, path string) error {

	// Append the line
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open kettle shell profile for writing: %w", err)
	}
	defer fileClose(file, &err)

	_, err = fmt.Fprintln(file, input)
	if err != nil {
		PrintErrors(err)
		return fmt.Errorf("could not open kettle shell profile for writing: %w", err)

	}
	return err

}

// AddLineToKettleShellProfile adds a given line of text to the kettle-specific shell profile.
func AddLineToKettleShellProfile(line string) error {
	shellInfo, err := GetShellInfo()
	if err != nil {
		return err
	}

	// Check if the line already exists
	shellExists, err := ExistsInFile(shellInfo.Path, line)
	if err != nil {
		return fmt.Errorf("error checking kettle shell profile: %w", err)
	}

	if shellExists {
		PrintInfo(fmt.Sprintf("Shell profile already contains %q.", line))
		return nil // Line already there, do nothing.
	}

	kettleExists, err := ExistsInFile(shellInfo.KettlePath, line)
	if err != nil {
		return fmt.Errorf("error checking kettle shell profile: %w", err)
	}

	if kettleExists {
		PrintInfo(fmt.Sprintf("%q already contains %q.", shellInfo.KettlePath, line))
		return nil // Line already there, do nothing.
	}

	err = addToFile(line, shellInfo.KettlePath)
	if err != nil {
		PrintErrors(err)
		return fmt.Errorf("could not add line to kettle shell profile: %w", err)
	}

	return nil
}

// EnsureCompletionsSourced adds the source command for completions to the kettle shell profile.
func EnsureCompletionsSourced() error {
	configDir, err := GetKettleConfigDir()
	if err != nil {
		return err
	}
	shell := GetCurrentShell()
	completionFile := filepath.Join(configDir, "completions", fmt.Sprintf("kettle.%s", shell))

	sourceCmd := fmt.Sprintf("source %s", completionFile)
	return AddLineToKettleShellProfile(sourceCmd)
}
