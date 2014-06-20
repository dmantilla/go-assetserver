package main

import (
	"./server"
)

// In development run like:
//
//   $ go build assetserver.go && FILES_ENV=environment ./assetserver
func main() {
	server.Run()
}
