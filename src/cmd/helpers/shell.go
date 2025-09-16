package helpers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
func GetShellProfile() (string, error) {
	shell := GetCurrentShell()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}

	switch shell {
	case "bash":
		return filepath.Join(homeDir, ".bashrc"), nil
	case "zsh":
		return filepath.Join(homeDir, ".zshrc"), nil
	case "fish":
		return filepath.Join(homeDir, ".config", "fish", "config.fish"), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

// AddLineToShellProfile adds a given line of text to the appropriate shell profile file
// if it does not already exist in the file.
func AddLineToShellProfile(line string) error {
	profilePath, err := GetShellProfile()
	if err != nil {
		return err
	}

	// Check if the line already exists in the file.
	exists, err := ExistsInFile(profilePath, line)
	if err != nil {
		return fmt.Errorf("error checking shell profile: %w", err)
	}
	if exists {
		PrintSuccess(fmt.Sprintf("Configuration already exists in %s.", filepath.Base(profilePath)))
		return nil
	}

	// If the file doesn't exist or the line isn't in it, append the line.
	file, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open shell profile for writing: %w", err)
	}
	defer file.Close()

	if _, err := fmt.Fprintln(file, line); err != nil {
		return fmt.Errorf("failed to write to shell profile: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Configuration added to %s.", filepath.Base(profilePath)))
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
	defer file.Close()

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

// GetKettleShellProfile returns the path to the kettle-specific shell profile file.
func GetKettleShellProfile() (string, error) {
	configDir, err := GetKettleConfigDir()
	if err != nil {
		return "", err
	}

	shell := GetCurrentShell()
	profileName := fmt.Sprintf("kettle.%src", shell) // e.g., kettle.zshrc

	return filepath.Join(configDir, profileName), nil
}

// EnsureKettleProfileSourced makes sure the main shell profile sources the kettle-specific profile.
func EnsureKettleProfileSourced() error {
	kettleProfile, err := GetKettleShellProfile()
	if err != nil {
		return err
	}

	sourceCmd := fmt.Sprintf("source %s", kettleProfile)

	return AddLineToShellProfile(sourceCmd)
}

// AddLineToKettleShellProfile adds a given line of text to the kettle-specific shell profile.
func AddLineToKettleShellProfile(line string) error {
	profilePath, err := GetKettleShellProfile()
	if err != nil {
		return err
	}

	// Check if the line already exists
	exists, err := ExistsInFile(profilePath, line)
	if err != nil {
		return fmt.Errorf("error checking kettle shell profile: %w", err)
	}
	if exists {
		PrintSuccess(fmt.Sprintf("Shell profile already contains %q.", line))
		return nil // Line already there, do nothing.
	}

	// Append the line
	file, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open kettle shell profile for writing: %w", err)
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, line)
	return err
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
