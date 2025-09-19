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

// CommandExists checks if a command is in the PATH or available as a shell function/builtin.
func CommandExists(cmd string) bool {
	shellInfo := GetShellInfo()
	// First try exec.LookPath for regular binaries
	_, err := exec.LookPath(cmd)
	if err == nil {
		PrintInfo(fmt.Sprintf("Command '%s' exists as binary.", cmd))
		return true
	}
	cmdStr := fmt.Sprintf("command -v %s", cmd)

	err = RunCmd(cmdStr, false)
	if err == nil {
		PrintInfo(fmt.Sprintf("Command '%s' exists", cmd))
		return true
	}

	// Use interactive shell to check for functions and builtins
	// The -i flag ensures we load shell functions like nvm
	cmdStr = fmt.Sprintf("%s -c 'source %s; %s'", shellInfo.ShellBinPath, shellInfo.ShellRCPath, cmdStr)
	err = RunCmd(cmdStr)
	if err != nil {
		return false
	}

	// PrintInfo(fmt.Sprintf("Command '%s' exists as shell function.", cmd))

	// // Check if output is not empty and not just whitespace
	// result := strings.TrimSpace(string(output))
	// if len(result) > 0 {
	// 	// Check if it's a function (functions typically contain newlines or "function" keyword)
	// 	if strings.Contains(result, "\n") || strings.Contains(result, cmd) {
	// 		PrintInfo(fmt.Sprintf("Command '%s' exists as shell function.", cmd))
	// 		return true
	// 	} else {
	// 		PrintInfo(fmt.Sprintf("Command '%s' exists (%s).", cmd, result))
	// 		return true
	// 	}
	// }

	PrintInfo(fmt.Sprintf("Command '%s' does not exist.", cmd))
	return false
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

func RunCmdWithShellProfile(command string) error {
	shellInfo := GetShellInfo()
	// Use interactive shell to ensure profile is loaded
	cmdStr := fmt.Sprintf("%s -c 'source %s; %s'", shellInfo.ShellBinPath, shellInfo.ShellRCPath, command)

	return RunCmd(cmdStr)
}
