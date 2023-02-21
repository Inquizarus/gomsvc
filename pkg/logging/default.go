package logging

import "os"

var DefaultLogger = NewPlainLogger(os.Stdout, "gotnb")
