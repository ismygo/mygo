package mygo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoveEmptyDirs(t *testing.T) {
	testPath := "testdata/dir"

	// case 1
	if err := os.RemoveAll(testPath); nil != err {
		t.Errorf("clear test empty dir [%s] failed: %s", testPath, err)
	}

	a := filepath.Join(testPath, "a")
	if err := os.MkdirAll(a, 0755); nil != err {
		t.Errorf("make dir [%s] failed: %s", testPath, err)
	}

	if err := File.RemoveEmptyDirs(testPath); nil != err {
		t.Errorf("remove empty dirs failed: %s", err)
	}

	if File.IsDir(a) || File.IsDir(testPath) {
		t.Errorf("empty dir [%s] exists", a)
	}

	// case 2
	if err := os.RemoveAll(testPath); nil != err {
		t.Errorf("clear test empty dir [%s] failed: %s", testPath, err)
	}

	if err := os.MkdirAll(a, 0755); nil != err {
		t.Errorf("make dir [%s] failed: %s", testPath, err)
	}
	test := filepath.Join(a, "test")
	if err := os.WriteFile(test, []byte(""), 0644); nil != err {
		t.Errorf("write file [%s] failed: %s", test, err)
	}

	if err := File.RemoveEmptyDirs(testPath); nil != err {
		t.Errorf("remove empty dirs failed: %s", err)
	}

	if !File.IsDir(a) || !File.IsDir(testPath) {
		t.Errorf("empty dir [%s] exists", a)
	}

	// case 3
	if err := os.RemoveAll(testPath); nil != err {
		t.Errorf("clear test empty dir [%s] failed: %s", testPath, err)
	}

	if err := os.MkdirAll(a, 0755); nil != err {
		t.Errorf("make dir [%s] failed: %s", testPath, err)
	}

	if err := File.RemoveEmptyDirs(testPath, "a"); nil != err {
		t.Errorf("remove empty dirs failed: %s", err)
	}

	if !File.IsDir(a) || !File.IsDir(testPath) {
		t.Errorf("empty dir [%s] exists", a)
	}

	if err := os.RemoveAll(testPath); nil != err {
		t.Errorf("clear test empty dir [%s] failed: %s", testPath, err)
	}
}

func TestWriteFileSaferByReader(t *testing.T) {
	writePath := "testdata/filewrite.go"
	defer os.RemoveAll(writePath)
	if err := File.WriteFileSaferByReader(writePath, strings.NewReader("test"), 0644); nil != err {
		t.Errorf("write file [%s] failed: %s", writePath, err)
	}
}

func TestWriteFileSafer(t *testing.T) {
	writePath := "testdata/filewrite.go"
	defer os.RemoveAll(writePath)

	if err := os.WriteFile(writePath, []byte("0"), 0644); nil != err {
		t.Fatalf("write file [%s] failed: %s", writePath, err)
	}

	info, err := os.Stat(writePath)
	if nil != err {
		t.Fatalf("stat file [%s] failed: %s", writePath, err)
	}

	if err = File.WriteFileSafer(writePath, []byte("test"), 0644); nil != err {
		t.Errorf("write file [%s] failed: %s", writePath, err)
	}

	info, err = os.Stat(writePath)
	if nil != err {
		t.Fatalf("stat file [%s] failed: %s", writePath, err)
	}
	modTime2 := info.ModTime()
	t.Logf("file mod time [%v]", modTime2)
}

func TestIsHidden(t *testing.T) {
	filename := "./file.go"
	isHidden := File.IsHidden(filename)
	if isHidden {
		t.Error("file [" + filename + "] is not hidden")
	}
}

func TestGetFileSize(t *testing.T) {
	filename := "testdata/README.md"
	if size := File.GetFileSize(filename); size == -1 {
		t.Error(fmt.Sprintf("GetFileSize Failed, size: %d bytes", size))
	} else {
		t.Log(fmt.Sprintf("size: %d bytes", size))
	}
}
