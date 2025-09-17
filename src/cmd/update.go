package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mpaxson/kettle/src/cmd/helpers"
	"github.com/spf13/cobra"
)

const (
	githubReleasesURL    = "https://github.com/mpaxson/kettle/releases/latest/download/kettle"
	githubAPIReleasesURL = "https://api.github.com/repos/mpaxson/kettle/releases/latest"
)

// GitHubRelease represents the structure of a GitHub release API response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update kettle to the latest version",
	Long: `Download and install the latest version of kettle from GitHub releases.
This command handles the case where the current binary might be in use
by using a temporary file and atomic replacement.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateKettle()
	},
}

func updateKettle() error {
	helpers.PrintInfo("Checking for latest kettle version...")

	// Get latest release information from GitHub API
	latestRelease, err := getLatestRelease()
	if err != nil {
		helpers.PrintError("Failed to fetch latest release information", err)
		return err
	}

	latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")

	// Get clean version without 'v' prefix
	currentVersion := strings.TrimPrefix(Version, "v")

	helpers.PrintInfo(fmt.Sprintf("Current version: v%s", currentVersion))
	helpers.PrintInfo(fmt.Sprintf("Latest version: v%s", latestVersion)) // Compare versions
	if compareVersions(currentVersion, latestVersion) >= 0 {
		helpers.PrintSuccess("You are already running the latest version!")
		return nil
	}

	helpers.PrintInfo("A newer version is available. Starting update...")

	// Get the current executable path
	currentExe, err := os.Executable()
	if err != nil {
		helpers.PrintError("Failed to get current executable path", err)
		return err
	}

	// Create a temporary file for the download
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "kettle-update")

	helpers.PrintInfo("Downloading latest kettle binary...")

	if err := helpers.DownloadFile(tmpFile, githubReleasesURL); err != nil {
		helpers.PrintError("Failed to download latest version", err)
		return err
	}

	// Make the temporary file executable
	if err := os.Chmod(tmpFile, 0755); err != nil {
		helpers.PrintError("Failed to make downloaded binary executable", err)
		err := os.Remove(tmpFile) // Clean up
		if err != nil {
			helpers.PrintError("Failed to remove temporary file", err)
			return err
		}
		return err
	}

	helpers.PrintInfo("Replacing current binary...")

	// Use atomic replacement to handle the case where the binary is in use
	if err := atomicReplace(currentExe, tmpFile); err != nil {
		helpers.PrintError("Failed to replace current binary", err)
		err := os.Remove(tmpFile) // Clean up
		if err != nil {
			helpers.PrintError("Failed to remove temporary file", err)
			return err
		}
		return err
	}

	helpers.PrintSuccess(fmt.Sprintf("Kettle has been successfully updated from v%s to v%s!", currentVersion, latestVersion))
	helpers.PrintInfo("You may need to restart your terminal session for changes to take effect.")

	return nil
}

// getLatestRelease fetches the latest release information from GitHub API
func getLatestRelease() (*GitHubRelease, error) {
	resp, err := http.Get(githubAPIReleasesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release information: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			helpers.PrintError("Failed to close response body", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release information: %w", err)
	}

	return &release, nil
}

// compareVersions compares two semantic version strings
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	// Handle dev versions
	if v1 == "dev" && v2 != "dev" {
		return -1 // dev is always older than any released version
	}
	if v1 != "dev" && v2 == "dev" {
		return 1 // any released version is newer than dev
	}
	if v1 == "dev" && v2 == "dev" {
		return 0 // dev == dev
	}

	// Parse version components
	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	// Compare each component
	for i := 0; i < 3; i++ {
		if parts1[i] < parts2[i] {
			return -1
		}
		if parts1[i] > parts2[i] {
			return 1
		}
	}

	return 0
}

// parseVersion parses a version string into [major, minor, patch] integers
func parseVersion(version string) [3]int {
	parts := [3]int{0, 0, 0}
	versionParts := strings.Split(version, ".")

	for i, part := range versionParts {
		if i >= 3 {
			break
		}
		// Simple integer parsing - ignoring errors and defaulting to 0
		if val := 0; len(part) > 0 {
			for _, char := range part {
				if char >= '0' && char <= '9' {
					val = val*10 + int(char-'0')
				} else {
					break // Stop at first non-digit
				}
			}
			parts[i] = val
		}
	}

	return parts
}

// atomicReplace performs an atomic replacement of the target file with the source file
// This handles the case where the binary might be in use on different operating systems
func atomicReplace(target, source string) error {
	switch runtime.GOOS {
	case "windows":
		// On Windows, we need to use a different approach since files can't be replaced while in use
		return windowsReplace(target, source)
	default:
		// On Unix-like systems, we can use rename which is atomic
		return unixReplace(target, source)
	}
}

// unixReplace performs atomic replacement on Unix-like systems
func unixReplace(target, source string) error {
	// On Unix systems, rename is atomic and can replace a file even if it's currently executing
	if err := os.Rename(source, target); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}
	return nil
}

// windowsReplace handles binary replacement on Windows systems
func windowsReplace(target, source string) error {
	// On Windows, we create a batch script that waits for the current process to exit
	// then replaces the binary
	batchScript := target + ".update.bat"

	batchContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
move "%s" "%s"
del "%%~f0"
`, source, target)

	if err := os.WriteFile(batchScript, []byte(batchContent), 0755); err != nil {
		return fmt.Errorf("failed to create update script: %w", err)
	}

	// Start the batch script in the background
	cmd := exec.Command("cmd", "/C", "start", "/B", batchScript)
	if err := cmd.Start(); err != nil {
		helpers.PrintError("Failed to start update script", err)
		return err
	}

	helpers.PrintInfo("Update script started. The binary will be replaced after this process exits.")
	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
