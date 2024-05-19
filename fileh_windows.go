//go:build windows

package mygo

import (
	"path/filepath"
	"syscall"
)

// IsHidden checks whether the file specified by the given path is hidden.
func (*GoFile) IsHidden(path string) bool {
	if baseName := filepath.Base(path); 1 <= len(baseName) && "." == baseName[:1] {
		return true
	}

	pointer, err := syscall.UTF16PtrFromString(path)
	if nil != err {
		logger.Errorf("Checks file [%s] is hidden failed: [%s]", path, err)
		return false
	}

	attributes, err := syscall.GetFileAttributes(pointer)
	if nil != err {
		logger.Errorf("Checks file [%s] is hidden failed: [%s]", path, err)
		return false
	}
	return 0 != attributes&syscall.FILE_ATTRIBUTE_HIDDEN
}
