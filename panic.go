package mygo

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
)

var (
	dunno     = []byte("???") // 未知错误
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recover 对 recover() 的封装，增加stack的信息记录、输出
func (*GoPanic) Recover(err *error) {
	if e := recover(); e != nil {
		stack := stack()
		msg := fmt.Sprintf("Panic Recover: %v\n\t%s\n", e, stack)
		logger.Errorf(msg)

		if err != nil {
			*err = errors.New(msg)
		}
	}
}

func stack() []byte {
	buf := &bytes.Buffer{} // 存储结果

	var (
		sourceCode [][]byte // 存储出现错误的代码文件
		lastFile   string   // 文件指针，用以作为标识
	)
	for i := 2; ; i++ {
		// file: 发送panic的文件，line：对应的代码行数，pc：16进制
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 输出堆栈
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)

		// 如果file == lastFile，说明已经读取到末尾，不再有更深层的堆栈信息
		if file != lastFile {
			data, err := os.ReadFile(file) // datas：读取file后的文件信息
			if err != nil {
				continue
			}

			// 读取到的文件信息
			sourceCode = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}

		// 代码行数修正。因为在array中以0开头，而在现实中以1开头
		line--

		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(sourceCode, line))
	}

	return buf.Bytes()
}

// source 返回出错文件对应的函数
func source(lines [][]byte, n int) []byte {
	// 未知错误
	if n < 0 || n >= len(lines) {
		return dunno
	}

	return bytes.Trim(lines[n], " \t")
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}

	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}

	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
