package helpers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
)

func GithubDownloadLatestRelease(owner, repo, destDir, binaryName string) (string, error) {
	ctx := context.Background()

	client := github.NewClient(nil)

	// Get latest release
	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return "", fmt.Errorf("failed to get latest release: %w", err)
	}

	// Collect all asset names and select the best one
	var assetNames []string
	for _, asset := range release.Assets {
		assetNames = append(assetNames, asset.GetName())
	}

	// Select the best asset using ranking
	bestAssetName := SelectBestAsset(assetNames)
	if bestAssetName == "" {
		return "", fmt.Errorf("no suitable asset found for %s/%s release %s", owner, repo, release.GetTagName())
	}

	// Find the download URL for the best asset
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.GetName() == bestAssetName {
			downloadURL = asset.GetBrowserDownloadURL()
			break
		}
	}

	assetName := bestAssetName

	// Download the file
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to download asset: %w", err)
	}
	defer IOClose(resp.Body, &err)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status downloading asset: %s", resp.Status)
	}

	// Save to destDir
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create dir %s: %w", destDir, err)
	}
	destPath := filepath.Join(destDir, assetName)

	out, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer IOClose(out, &err)

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save asset: %w", err)
	}

	PrintSuccess(destPath + " downloaded successfully")

	// Check if the downloaded file is an archive and extract if needed
	if isArchive(assetName) {
		PrintInfo("Extracting binary from archive...")

		// Extract the binary from the archive
		if err := ExtractBinaryFromArchive(destPath, destDir, binaryName); err != nil {
			return "", fmt.Errorf("failed to extract binary from archive: %w", err)
		}

		// Remove the archive file after extraction
		if err := os.Remove(destPath); err != nil {
			PrintError("Failed to remove archive file", err)
		}

		// Return the path to the extracted binary
		finalBinaryPath := filepath.Join(destDir, binaryName)
		PrintSuccess(fmt.Sprintf("Binary extracted to: %s", finalBinaryPath))
		return finalBinaryPath, nil
	}

	// If it's already a binary, rename it to the expected binary name
	finalBinaryPath := filepath.Join(destDir, binaryName)
	if destPath != finalBinaryPath {
		if err := os.Rename(destPath, finalBinaryPath); err != nil {
			return "", fmt.Errorf("failed to rename binary: %w", err)
		}
		if err := os.Chmod(finalBinaryPath, 0755); err != nil {
			return "", fmt.Errorf("failed to make binary executable: %w", err)
		}
	}

	return finalBinaryPath, nil
}

// GithubGetLatestRelease gets the latest release information from a GitHub repository
func GithubGetLatestRelease(owner, repo string) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release for %s/%s: %w", owner, repo, err)
	}

	return release, nil
}
