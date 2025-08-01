package helper

import (
	"io"
	"os"
)

// credits to https://stackoverflow.com/a/29339052
func CaptureStdout(f func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()

	os.Stdout = w
	defer func() {
		os.Stdout = rescueStdout
	}()

	f()

	w.Close()
	b, _ := io.ReadAll(r)

	return string(b)
}
