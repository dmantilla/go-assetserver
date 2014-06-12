package main

import (
	"./server"
)

// In development run like:
//
//   $ go build imageserver.go && FILES_ENV=environment ./imageserver
func main() {
	server.Run()
}
