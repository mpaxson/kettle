package tests

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/mpaxson/kettle/src/cmd/helpers"
	"github.com/stretchr/testify/assert"
)

func TestAssetSelection(t *testing.T) {
	assets := []string{
		"source.tar.gz",
		"golangci-lint-1.54.2-linux-amd64.tar.gz",
		"golangci-lint-1.54.2-linux-amd64.deb",
		"golangci-lint-1.54.2-linux-amd64",
		"golangci-lint-1.54.2-darwin-amd64.tar.gz",
		"v1.54.2.tar.gz",
		"golangci-lint-1.54.2-windows-amd64.zip",
	}
	fmt.Println("Asset rankings:")
	for _, asset := range assets {
		rank := helpers.RankAsset(asset)
		fmt.Printf("  %s: %d\n", asset, rank)
	}

	best := helpers.SelectBestAsset(assets)
	if runtime.GOARCH == "amd64" && runtime.GOOS == "linux" {

		assert.Equal(t, "golangci-lint-1.54.2-linux-amd64", best)
	} else if runtime.GOARCH == "amd64" && runtime.GOOS == "darwin" {
		assert.Equal(t, "golangci-lint-1.54.2-darwin-amd64.tar.gz", best)

	}

}
