package helpers

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// DownloadFile downloads a file from the given URL and saves it to the specified path

func DownloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer IOClose(out, &err)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			PrintError("Failed to close response body", err)
		}
	}()
	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	err = os.Chmod(filepath, 0755)
	if err != nil {
		return fmt.Errorf("failed to make file executable: %w", err)
	}
	return nil
}

// AssetRank represents a ranked asset
type AssetRank struct {
	Name string
	Rank int
}

// MatchAsset picks the best asset for the current platform
func MatchAsset(name string) bool {
	return RankAsset(name) > 0
}

// RankAsset assigns a rank to an asset for the current platform
// Higher rank = better match. 0 = not suitable
func RankAsset(name string) int {
	name = strings.ToLower(name)
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH

	// Ignore source archives completely
	if isSourceArchive(name) {
		return 0
	}

	// Convert GOARCH to common naming conventions
	arch := currentArch
	if currentArch == "amd64" {
		arch = "amd64" // Also matches x86_64, x64
	}

	// Check if name contains OS and architecture
	hasOS := strings.Contains(name, currentOS)
	hasArch := strings.Contains(name, arch) ||
		strings.Contains(name, "x86_64") ||
		strings.Contains(name, "x64")

	// Must match architecture to be considered
	if !hasArch {
		return 0
	}

	// Base ranking system
	rank := 0

	// Platform compatibility
	switch currentOS {
	case "linux":
		if hasOS || strings.HasSuffix(name, ".deb") || strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".zip") {
			rank += 10
		}
	case "darwin":
		if hasOS && (strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".dmg") || strings.HasSuffix(name, ".zip")) {
			rank += 10
		}
	case "windows":
		if hasOS && (strings.HasSuffix(name, ".zip") || strings.HasSuffix(name, ".exe")) {
			rank += 10
		}
	}

	if rank == 0 {
		return 0 // Not compatible with current platform
	}

	// Priority ranking: standalone-binary > archive > deb
	if isStandaloneBinary(name) {
		rank += 100
	} else if isArchive(name) {
		rank += 50
	} else if isDeb(name) {
		rank += 25
	}

	// Bonus for exact OS match
	if hasOS {
		rank += 5
	}

	// Bonus for Ubuntu .deb packages on Ubuntu systems
	if currentOS == "linux" && IsUbuntu() && strings.HasSuffix(name, ".deb") {
		rank += 15
	}

	return rank
}

// SelectBestAsset selects the best asset from a list of asset names
func SelectBestAsset(assetNames []string) string {
	var bestAsset AssetRank
	bestAsset.Rank = 0

	for _, name := range assetNames {
		rank := RankAsset(name)
		if rank > bestAsset.Rank {
			bestAsset.Name = name
			bestAsset.Rank = rank
		}
	}

	return bestAsset.Name
}

// isSourceArchive checks if the asset is a source code archive
func isSourceArchive(name string) bool {
	return strings.Contains(name, "source") ||
		strings.Contains(name, "src") ||
		name == "source.tar.gz" ||
		name == "source.zip" ||
		strings.HasPrefix(name, "v") && (strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".zip")) && !strings.Contains(name, runtime.GOOS) && !strings.Contains(name, runtime.GOARCH)
}

// isStandaloneBinary checks if the asset is a standalone binary
func isStandaloneBinary(name string) bool {
	return isExecutable(name) || (!isArchive(name) && !isPackage(name))

}

func isExecutable(name string) bool {
	return !strings.Contains(name, ".") || strings.HasSuffix(name, ".exe") || strings.HasSuffix(name, runtime.GOOS+"-"+runtime.GOARCH)
}

// isArchive checks if the asset is an archive (tar.gz, zip, etc.)
func isArchive(name string) bool {
	return strings.HasSuffix(name, ".tar.gz") ||
		strings.HasSuffix(name, ".zip") ||
		strings.HasSuffix(name, ".tar") ||
		strings.HasSuffix(name, ".gz") ||
		strings.HasSuffix(name, ".dmg")
}
func isPackage(name string) bool {
	return isDeb(name) ||
		isRPM(name)
}

// isDeb checks if the asset is a .deb package
func isDeb(name string) bool {
	return strings.HasSuffix(name, ".deb")
}
func isRPM(name string) bool {
	return strings.HasSuffix(name, ".rpm")
}

// ExtractBinaryFromArchive extracts a binary from an archive and places it in destDir
func ExtractBinaryFromArchive(archivePath, destDir, binaryName string) error {
	ext := strings.ToLower(filepath.Ext(archivePath))

	switch {
	case strings.HasSuffix(archivePath, ".tar.gz"):
		return extractFromTarGz(archivePath, destDir, binaryName)
	case ext == ".zip":
		return extractFromZip(archivePath, destDir, binaryName)
	case ext == ".deb":
		// For .deb files, extract the binary from the data.tar.* inside
		return extractFromDeb(archivePath, destDir, binaryName)
	default:
		// Assume it's already a binary, just move it
		finalPath := filepath.Join(destDir, binaryName)
		if err := os.Rename(archivePath, finalPath); err != nil {
			return fmt.Errorf("failed to move binary: %w", err)
		}
		if err := os.Chmod(finalPath, 0755); err != nil {
			return fmt.Errorf("failed to make binary executable: %w", err)
		}
		return nil
	}
}

// extractFromTarGz extracts a binary from a .tar.gz archive
func extractFromTarGz(archivePath, destDir, binaryName string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer IOClose(file, &err)

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer IOClose(gzr, &err)

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Look for the binary (could be in subdirectories)
		if strings.HasSuffix(header.Name, binaryName) && header.Typeflag == tar.TypeReg {
			destPath := filepath.Join(destDir, binaryName)
			outFile, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				IOClose(outFile, &err)
				return fmt.Errorf("failed to extract binary: %w", err)
			}
			defer IOClose(outFile, &err)

			if err := os.Chmod(destPath, 0755); err != nil {
				return fmt.Errorf("failed to make binary executable: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("binary %s not found in archive", binaryName)
}

// extractFromZip extracts a binary from a .zip archive
func extractFromZip(archivePath, destDir, binaryName string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer IOClose(r, &err)

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, binaryName) && !f.FileInfo().IsDir() {
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open file in zip: %w", err)
			}

			destPath := filepath.Join(destDir, binaryName)
			outFile, err := os.Create(destPath)
			if err != nil {
				err2 := rc.Close()
				if err2 != nil {
					PrintError("Failed to close zip file", err2)
					return err2
				}
				return fmt.Errorf("failed to create output file: %w", err)
			}

			if _, err := io.Copy(outFile, rc); err != nil {
				IOClose(outFile, &err)
				IOClose(rc, &err)
				return fmt.Errorf("failed to extract binary: %w", err)
			}

			IOClose(outFile, &err)
			IOClose(rc, &err)

			if err := os.Chmod(destPath, 0755); err != nil {
				return fmt.Errorf("failed to make binary executable: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("binary %s not found in zip archive", binaryName)
}

// extractFromDeb extracts a binary from a .deb package
func extractFromDeb(archivePath, destDir, binaryName string) error {
	// For now, this is a simplified implementation
	// In a full implementation, you'd extract the data.tar.* from the .deb
	// and then extract the binary from that
	return fmt.Errorf("deb extraction not yet implemented - please use tar.gz version")
}

// Download and install a script from a url
func DownloadAndRunInstallScript(url string, filename string) error {
	curDir := GetCurrentDir()

	shScriptPath := filepath.Join(curDir, filename)
	if err != nil {
		PrintError("Failed to get latest NVM release info", err)
		return err
	}

	if err := DownloadFile(shScriptPath, url); err != nil {
		PrintError("Failed to download install script", err)
		return err
	}
	const halfSecond = 500 * time.Millisecond
	time.Sleep(halfSecond)

	cmd := fmt.Sprintf("cat %s | bash", shScriptPath)
	if err := RunCmd(cmd); err != nil {
		PrintError("Failed to run install script", err)
		return err
	}
	err := os.Remove(shScriptPath)
	if err != nil {
		PrintError("Failed to remove temporary script", err)
		return err
	}
	return nil
}
