package mygo

import (
	"runtime"
	"testing"
)

func TestIsWindows(t *testing.T) {
	goos := runtime.GOOS

	if "windows" == goos && !OS.IsWindows() {
		t.Error("runtime.GOOS returns [windows]")
		return
	}
}

func TestIsLinux(t *testing.T) {
	goos := runtime.GOOS

	if "linux" == goos && !OS.IsLinux() {
		t.Error("runtime.GOOS returns [linux]")
		return
	}
}

func TestIsDarwin(t *testing.T) {
	goos := runtime.GOOS

	if "darwin" == goos && !OS.IsDarwin() {
		t.Error("runtime.GOOS returns [darwin]")
		return
	}
}

func TestPwd(t *testing.T) {
	pwd := OS.Pwd()
	t.Log(pwd)
}

func TestHome(t *testing.T) {
	home, err := OS.Home()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(home)
}
