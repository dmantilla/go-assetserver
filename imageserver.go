package main

import (
	"./server"
)

// In development run like:
//
//   $ go build main.go && FILES_ENV=environment ./main
func main() {
	server.Run()
}
