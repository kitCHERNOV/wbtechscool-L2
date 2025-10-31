package main

// wget util on Go

import (
	"wget/include/wgetmanager"
	"wget/pkg/directories"
)

func main() {
	// start service
	wgetmanager.WgetManager()
	directories.ToCreateDirectories()
}
