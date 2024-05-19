package mygo

import "os"

// Logger for mygo
var logger = Log.NewLogger(os.Stdout)

type (
	GoNet   byte
	GoOS    byte
	GoPanic byte
)

var (
	Net   GoNet
	OS    GoOS
	Panic GoPanic
)
