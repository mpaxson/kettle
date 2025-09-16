package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mpaxson/kettle/src/cmd/helpers"
	"github.com/spf13/cobra"
)

const (
	githubReleasesURL = "https://github.com/mpaxson/kettle/releases/latest/download/kettle"
)

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

	if err := downloadFile(tmpFile, githubReleasesURL); err != nil {
		helpers.PrintError("Failed to download latest version", err)
		return err
	}

	// Make the temporary file executable
	if err := os.Chmod(tmpFile, 0755); err != nil {
		helpers.PrintError("Failed to make downloaded binary executable", err)
		os.Remove(tmpFile) // Clean up
		return err
	}

	helpers.PrintInfo("Replacing current binary...")

	// Use atomic replacement to handle the case where the binary is in use
	if err := atomicReplace(currentExe, tmpFile); err != nil {
		helpers.PrintError("Failed to replace current binary", err)
		os.Remove(tmpFile) // Clean up
		return err
	}

	helpers.PrintSuccess("Kettle has been successfully updated to the latest version!")
	helpers.PrintInfo("You may need to restart your terminal session for changes to take effect.")

	return nil
}

// downloadFile downloads a file from the given URL and saves it to the specified path
func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer helpers.FileClose(out, &err)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
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
	cmd.Start()

	helpers.PrintInfo("Update script started. The binary will be replaced after this process exits.")
	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
