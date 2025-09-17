// Package helpers provides utility functions for OS detection and command execution.
package helpers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// getOSReleaseValue reads a specific key from the /etc/os-release file.

type OSRelease map[string]string

// IOClose ensures file.Close() errors don't get lost.
// If *err is nil, it will overwrite it with the close error.
// If *err is already set, it will preserve the original error.
func IOClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil {
		if *err == nil {
			*err = cerr
		}
	}
}

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
	defer IOClose(file, &err)

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

// InstallBinary copies the current executable to a directory in the user's PATH.
func InstallBinary(path string) {
	// 1. Get the path of the current executable

	// 2. Determine the installation directory
	destDir, err := GetInstallDir()
	if err != nil {
		PrintError("", err)
	}

	destPath := filepath.Join(destDir, filepath.Base(path))

	// 3. Open the source file
	srcFile, err := os.Open(path)
	if err != nil {
		PrintErrors(err)
	}
	defer IOClose(srcFile, &err)

	// 4. Create the destination file (overwrite if it exists)
	destFile, err := os.Create(destPath)
	if err != nil {
		PrintError("could not create destination file:", err)
		return
	}
	defer IOClose(destFile, &err)

	// 5. Copy the file contents
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		PrintError("could not copy binary:", err)
		return
	}

	// 6. Make the destination file executable
	err = os.Chmod(destPath, 0755)
	if err != nil {
		PrintError("Could not make binary executable:", err)
		return
	}

	if !CommandExists(filepath.Base(path)) {
		PrintError(fmt.Sprintf("Install Failed for %s", filepath.Base(path)))
		return
	}
	PrintSuccess(fmt.Sprintf("%s Installed Successfully", filepath.Base(path)))
}

// GetInstallDir determines the correct directory for installation.
func GetInstallDir() (string, error) {
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
