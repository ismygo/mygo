package mygo

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// IsWindows determines whether current OS is Windows.
func (*GoOS) IsWindows() bool {
	return "windows" == runtime.GOOS
}

// IsLinux determines whether current OS is Linux.
func (*GoOS) IsLinux() bool {
	return "linux" == runtime.GOOS
}

// IsDarwin determines whether current OS is Darwin.
func (*GoOS) IsDarwin() bool {
	return "darwin" == runtime.GOOS
}

// Pwd 获取当前工作路径的地址
func (*GoOS) Pwd() string {
	file, _ := exec.LookPath(os.Args[0])
	pwd, _ := filepath.Abs(file)

	return filepath.Dir(pwd)
}

// Home 获取当前家目录。
// Home 的实现在不同OS有不同的实现，获取方法取决于不同OS
func (*GoOS) Home() (string, error) {
	user, err := user.Current()
	if err == nil {
		return user.HomeDir, nil
	}

	// Support Windows
	if OS.IsWindows() {
		return homeWindows()
	}

	// Support Unix-like OS
	return homeUnix()
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

func homeUnix() (string, error) {
	//  首先尝试直接获取HOME环境变量
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// 如果失败，则尝试shell脚本
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}
