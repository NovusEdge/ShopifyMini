package main

import (
	cp "cmputil/cmputils"
)

func main() {
	f := cp.File{Path: "./tempFile.js"}
	f.Gzip("./")
	f.Tar("./")
}
