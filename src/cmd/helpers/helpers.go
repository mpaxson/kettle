// Package helpers provides utility functions for OS detection and command execution.
package helpers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// getOSReleaseValue reads a specific key from the /etc/os-release file.

type OSRelease map[string]string

var (
	cached OSRelease
	once   sync.Once
	err    error
)

// Get returns the OSRelease map, reading /etc/os-release once
func Get() (OSRelease, error) {
	once.Do(func() {
		cached, err = readOSRelease("/etc/os-release")
	})
	return cached, err
}

func readOSRelease(path string) (OSRelease, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make(OSRelease)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		val := strings.Trim(parts[1], `"`)
		data[key] = val
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func IsUbuntu() bool {
	osr, err := Get()
	if err != nil {
		return false
	}
	return osr["ID"] == "ubuntu"
}

func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// CommandExists checks if a command is in the PATH.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	val := err == nil
	if val {
		PrintSuccess(successStyle.Render("Command '%s' exists.", cmd))
	}
	return val
}

func isUbuntuVersion(versionPrefix string) bool {
	if IsUbuntu() {
		osr, err := Get()
		if err != nil {
			return false
		}
		versionID := osr["VERSION_ID"]
		return strings.HasPrefix(versionID, versionPrefix)
	}
	return false
}

func IsUbuntu22() bool {
	return isUbuntuVersion("22.04")
}

func IsUbuntu24() bool {
	return isUbuntuVersion("24.04")
}

func IsUbuntu26() bool {
	return isUbuntuVersion("26.04")
}

// InstallBinary copies the current executable to a directory in the user's PATH.
func InstallBinary(path string) error {
	// 1. Get the path of the current executable

	// 2. Determine the installation directory
	destDir, err := getInstallDir()
	if err != nil {
		PrintError("", err)
		return err
	}

	destPath := filepath.Join(destDir, filepath.Base(path))

	// 3. Open the source file
	srcFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open source binary: %w", err)
	}
	defer srcFile.Close()

	// 4. Create the destination file (overwrite if it exists)
	destFile, err := os.Create(destPath)
	if err != nil {
		PrintError("could not create destination file:", err)
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer destFile.Close()

	// 5. Copy the file contents
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		PrintError("could not copy binary:", err)
		return fmt.Errorf("could not copy binary: %w", err)
	}

	// 6. Make the destination file executable
	err = os.Chmod(destPath, 0755)
	if err != nil {
		PrintError("Could not make binary executable:", err)
		return fmt.Errorf("could not make binary executable: %w", err)
	}

	return nil
}

// getInstallDir determines the correct directory for installation.
func getInstallDir() (string, error) {
	// Check for ~/.local/bin
	homeDir, err := os.UserHomeDir()
	if err != nil {
		PrintError("could not get user home directory:", err)
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}
	localBin := filepath.Join(homeDir, ".local", "bin")

	// Check if the directory exists
	info, err := os.Stat(localBin)
	if err == nil && info.IsDir() {
		// Check if it's in the PATH
		path := os.Getenv("PATH")
		if strings.Contains(path, localBin) {
			return localBin, nil
		}
	}

	// Fallback to /usr/local/bin
	return "/usr/local/bin", nil
}
