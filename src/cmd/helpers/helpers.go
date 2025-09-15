package helpers

import (
	"bufio"
	"os"
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
