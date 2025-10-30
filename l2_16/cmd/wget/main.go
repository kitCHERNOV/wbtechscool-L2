package main

// wget util on Go

import (
	"wget/include/wgetmanager"
	"wget/pkg/directories"
)

func main() {
	// start service
	// TODO: add getting urls ant throw down to WgetManager
	wgetmanager.WgetManager()
	directories.ToCreateDirectories()
}
