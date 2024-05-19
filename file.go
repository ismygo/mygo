package mygo

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RemoveEmptyDirs 移除指定文件夹下的空文件夹
func (GoFile) RemoveEmptyDirs(dir string, excludes ...string) (err error) {
	_, err = removeEmptyDirs(dir, excludes...)
	return
}

// removeEmptyDirs 递归移除文件夹下全部空文件夹
func removeEmptyDirs(dir string, excludes ...string) (removed bool, err error) {
	dirName := filepath.Base(dir)
	if Str.Contains(dirName, excludes) {
		return
	}

	var hasEntries bool             // 空文件夹标识
	entires, err := os.ReadDir(dir) // 读取该目录下的所有文件
	if err != nil {
		return false, err
	}
	for _, entry := range entires {
		if entry.IsDir() {
			subDir := filepath.Join(dir, entry.Name())
			// 递归清除
			removed, err = removeEmptyDirs(subDir, excludes...)
			if err != nil {
				return false, err
			}
			if !removed {
				hasEntries = true
			}
		} else {
			hasEntries = true
		}
	}

	if !hasEntries {
		if err = os.Remove(dir); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// IsDir 判断是否是一个文件夹
func (*GoFile) IsDir(path string) bool {
	file, err := os.Lstat(path)
	if os.IsNotExist(err) { // 目录不存在
		return false
	}
	if err != nil { // 其他错误
		logger.Warnf("determines whether [%s] is a directory failed: [%v]", path, err)
		return false
	}
	return file.IsDir()
}

// IsValidFileName 非法名字检查
func (GoFile) IsValidFileName(name string) bool {
	reserved := []string{"\\", "/", ":", "*", "?", "\"", "'", "<", ">", "|"}
	for _, r := range reserved {
		if strings.Contains(name, r) {
			return false
		}
	}
	return true
}

// WriteFileSaferByReader 方法通过提供的 reader 安全地写入文件到指定路径，使用指定的文件权限。
func (GoFile) WriteFileSaferByReader(writePath string, reader io.Reader, perm os.FileMode) (err error) {
	// 分割路径为目录和文件名
	dir, name := filepath.Split(writePath)
	// 创建临时文件路径
	tmp := filepath.Join(dir, name+Rand.String(7)+".tmp")
	// 尝试打开或创建临时文件，仅当文件不存在时创建
	f, err := os.OpenFile(tmp, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0660)
	if err != nil {
		return
	}

	// 从reader中复制数据到文件
	if _, err = io.Copy(f, reader); err != nil {
		return
	}

	// 同步文件内容到磁盘
	if err = f.Sync(); err != nil {
		return
	}

	// 关闭文件
	if err = f.Close(); err != nil {
		return
	}

	// 更改文件权限
	if err = os.Chmod(f.Name(), perm); err != nil {
		return
	}

	// 尝试重命名文件，最多重试3次
	for i := 0; i < 3; i++ {
		// 在Windows上，重命名操作不是原子的
		if err = os.Rename(f.Name(), writePath); err == nil {
			os.Remove(f.Name()) // 删除临时文件
			return
		}

		// 如果错误是因为访问被拒绝或文件被其他进程使用，则等待200毫秒后重试
		if errMsg := strings.ToLower(err.Error()); strings.Contains(errMsg, "access is denied") || strings.Contains(errMsg, "used by another process") {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		break
	}
	return
}

// WriteFileSafer 文件写入
func (GoFile) WriteFileSafer(writePath string, data []byte, perm os.FileMode) (err error) {
	dir, name := filepath.Split(writePath)

	// 创建临时文件
	tmp := filepath.Join(dir, name+Rand.String(7)+".tmp")
	f, err := os.OpenFile(tmp, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return
	}

	// 读取到临时文件
	if _, err = f.Write(data); err != nil {
		return
	}

	// 磁盘刷盘
	if err = f.Sync(); err != nil {
		return
	}

	// 关闭文件流
	if err = f.Close(); err != nil {
		return
	}

	// 修改文件权限
	if err = os.Chmod(f.Name(), perm); err != nil {
		return
	}

	// 重命名，默认重试 3 次
	for i := 0; i < 3; i++ {
		if err = os.Rename(f.Name(), writePath); err == nil {
			os.Remove(f.Name())
			return
		}

		if errMsg := strings.ToLower(err.Error()); strings.Contains(errMsg, "access is denied") || strings.Contains(errMsg, "used by another process") { // 文件可能是被锁定
			time.Sleep(200 * time.Millisecond)
			continue
		}
		break
	}
	return
}

// GetFileSize 获取文件大小(字节大小)
func (*GoFile) GetFileSize(path string) int64 {
	file, err := os.Stat(path)
	if err != nil {
		logger.Error(err)
		return -1
	}
	return file.Size()
}

// IsExist 文件是否存在
func (*GoFile) IsExist(path string) bool {
	_, err := os.Stat(path)

	return err == nil || os.IsExist(err)
}

// IsImg determines whether the specified extension is a image.
func (*GoFile) IsImg(extension string) bool {
	ext := strings.ToLower(extension)

	switch ext {
	case ".jpg", ".jpeg", ".bmp", ".gif", ".png", ".svg", ".ico":
		return true
	default:
		return false
	}
}

// Copy copies the source to the dest.
// Keep the dest access/mod time as the same as the source.
func (gl *GoFile) Copy(source, dest string) (err error) {
	if !gl.IsExist(source) {
		return os.ErrNotExist
	}

	if gl.IsDir(source) {
		return gl.copyDir(source, dest, false, true)
	}
	return gl.copyFile(source, dest, false, true)
}

func (gl *GoFile) copyDir(source, dest string, ignoreHidden, chtimes bool) (err error) {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if ignoreHidden && gl.IsHidden(source) {
		return
	}

	if err = os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	dirs, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, f := range dirs {
		srcFilePath := filepath.Join(source, f.Name())
		destFilePath := filepath.Join(dest, f.Name())

		if f.IsDir() {
			err = gl.copyDir(srcFilePath, destFilePath, ignoreHidden, chtimes)
			if err != nil {
				logger.Error(err)
				return
			}
		} else {
			err = gl.copyFile(srcFilePath, destFilePath, ignoreHidden, chtimes)
			if err != nil {
				logger.Error(err)
				return
			}
		}
	}

	if chtimes {
		if err = os.Chtimes(dest, sourceInfo.ModTime(), sourceInfo.ModTime()); nil != err {
			return
		}
	}
	return nil
}

func (gl *GoFile) copyFile(source, dest string, ignoreHidden, chtimes bool) (err error) {
	sourceinfo, err := os.Lstat(source)
	if nil != err {
		return
	}

	if 0 != sourceinfo.Mode()&os.ModeSymlink {
		// 忽略符号链接
		return
	}

	if ignoreHidden && gl.IsHidden(source) {
		return
	}

	if err = os.MkdirAll(filepath.Dir(dest), 0755); nil != err {
		return
	}

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	if err = os.Chmod(dest, sourceinfo.Mode()); nil != err {
		return
	}

	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourcefile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if nil != err {
		return
	}

	if chtimes {
		if err = os.Chtimes(dest, sourceinfo.ModTime(), sourceinfo.ModTime()); nil != err {
			return
		}
	}
	return
}

// CopyWithoutHidden copies the source to the dest without hidden files.
func (gl *GoFile) CopyWithoutHidden(source, dest string) (err error) {
	if !gl.IsExist(source) {
		return os.ErrNotExist
	}

	if gl.IsDir(source) {
		return gl.copyDir(source, dest, true, true)
	}
	return gl.copyFile(source, dest, true, true)
}

// CopyNewtimes copies the source to the dest.
// Do not keep the dest access/mod time as the same as the source.
func (gl *GoFile) CopyNewtimes(source, dest string) (err error) {
	if !gl.IsExist(source) {
		return os.ErrNotExist
	}

	if gl.IsDir(source) {
		return gl.CopyDirNewtimes(source, dest)
	}
	return gl.CopyFileNewtimes(source, dest)
}

// CopyDirNewtimes copies the source directory to the dest directory.
// Do not keep the dest access/mod time as the same as the source.
func (gl *GoFile) CopyDirNewtimes(source, dest string) (err error) {
	return gl.copyDir(source, dest, false, false)
}

// CopyFileNewtimes copies the source file to the dest file.
// Do not keep the dest access/mod time as the same as the source.
func (gl *GoFile) CopyFileNewtimes(source, dest string) (err error) {
	return gl.copyFile(source, dest, false, false)
}

// CopyFile copies the source file to the dest file.
// Keep the dest access/mod time as the same as the source.
func (gl *GoFile) CopyFile(source, dest string) (err error) {
	return gl.copyFile(source, dest, false, true)
}

// CopyDir copies the source directory to the dest directory.
// Keep the dest access/mod time as the same as the source.
func (gl *GoFile) CopyDir(source, dest string) (err error) {
	return gl.copyDir(source, dest, false, true)
}
