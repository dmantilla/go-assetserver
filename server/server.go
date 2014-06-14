package server

import (
	"net/http"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"log"
	"github.com/gographics/imagick/imagick"
	"../config"
	"../provider"
)

var cfg config.Configuration
var amazon provider.Amazon
var logger *log.Logger

func handler(response http.ResponseWriter, request *http.Request) {
	defer timeTrack(time.Now(), request.Method + " " + request.URL.Path)
	amazon.WriteAsset(cfg.AwsNode("assets_bucket"), AssetName(request), response, request)
}

func createLog() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

func loadConfiguration() {
	var err error

	dir := currentDir()
	cfg, err = config.Load(dir)
	if err != nil {
		fmt.Printf("There was a problem loading %s in %s: %s\n", config.FileName(), dir, err)
		os.Exit(1)
	}
}

func connectToAmazon() {
	var err error

	amazon, err = provider.AWSConnect(cfg.AwsNode("access_key"), cfg.AwsNode("secret_key"), cfg.CacheNode("folder"))
	if err != nil {
		fmt.Println("Cannot connect to Amazon")
		os.Exit(1)
	}
}

func timeTrack(start time.Time, message string) {
	elapsed := time.Since(start)
	logger.Printf("%s took %s", message, elapsed)
}

func Run() {
	imagick.Initialize()
	defer imagick.Terminate()

	createLog()
	loadConfiguration()
	connectToAmazon()

	http.HandleFunc("/", handler)
	fmt.Println("Running...")
	http.ListenAndServe(cfg.Address, nil)
}

func AssetName(r *http.Request) string {
	return r.URL.Path[1:]
}

func currentDir() (dir string) {
	dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	return
}
