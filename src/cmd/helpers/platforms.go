package helpers

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	isDarwinOnce   sync.Once
	isDarwin       bool
	isUbuntuOnce   sync.Once
	isUbuntu       bool
	isUbuntu22Once sync.Once
	isUbuntu22     bool
	isUbuntu24Once sync.Once
	isUbuntu24     bool
	isUbuntu26Once sync.Once
	isUbuntu26     bool
)

func IsUbuntu() bool {
	isUbuntuOnce.Do(func() {
		osr, err := Get()
		if err != nil {
			isUbuntu = false
			return
		}
		isUbuntu = osr["ID"] == "ubuntu"
	})
	return isUbuntu
}

func IsDarwin() bool {
	isDarwinOnce.Do(func() {
		isDarwin = runtime.GOOS == "darwin"
	})
	return isDarwin
}

// CommandExists checks if a command is in the PATH.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	val := err == nil
	if val {
		PrintInfo(fmt.Sprintf("Command '%s' exists.", cmd))
	}
	return val
}

func IsUbuntuVersion(versionPrefix string) bool {
	if !IsUbuntu() {
		return false
	}

	osr, err := Get()
	if err != nil {
		return false
	}
	versionID := osr["VERSION_ID"]
	return strings.HasPrefix(versionID, versionPrefix)

}

func IsUbuntu22() bool {
	isUbuntu22Once.Do(func() {
		isUbuntu22 = IsUbuntuVersion("22.04")
	})
	return isUbuntu22
}

func IsUbuntu24() bool {
	isUbuntu24Once.Do(func() {
		isUbuntu24 = IsUbuntuVersion("24.04")
	})
	return isUbuntu24
}

func IsUbuntu26() bool {
	isUbuntu26Once.Do(func() {
		isUbuntu26 = IsUbuntuVersion("26.04")
	})
	return isUbuntu26
}
