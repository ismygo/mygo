//go:build !windows
// +build !windows

package mygo

import "path/filepath"

// IsHidden 检查是否被隐藏
func (*GoFile) IsHidden(path string) bool {
	path = filepath.Base(path)
	if 1 > len(path) {
		return false
	}
	return "." == path[:1]
}
