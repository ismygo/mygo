package mygo

import "os"

// Logger for mygo
var logger = Log.NewLogger(os.Stdout)

type (
	GoNet   byte
	GoOS    byte
	GoPanic byte
	GoJSON  byte
	GoFile  byte
	GoStr   byte
	GoRand  byte
)

var (
	Net   GoNet
	OS    GoOS
	Panic GoPanic
	JSON  GoJSON
	File  GoFile
	Str   GoStr
	Rand  GoRand
)
