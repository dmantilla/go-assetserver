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
	"../asset"
)

var cfg config.Configuration
var amazon provider.Amazon
var cache provider.CacheProvider
var logger *log.Logger

func handler(response http.ResponseWriter, request *http.Request) {
	defer timeTrack(time.Now(), request.Method + " " + request.URL.Path + "&" + request.URL.RawQuery)
	err := asset.WriteResponse(request, response, &amazon, &cache, logger)
	if err != nil {
		logger.Printf("ERROR: %s", err)
		http.Error(response, "Resource not found", 404)
	}
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

	amazon, err = provider.AWSConnect(cfg, logger)
	if err != nil {
		fmt.Println("Cannot connect to Amazon")
		os.Exit(1)
	}
}

func connectToCache() {
	cache = provider.CacheConnect(cfg)
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
	connectToCache()

	http.HandleFunc("/", handler)
	fmt.Printf("Running on %s... v0.9.0\n", cfg.Address)
	http.ListenAndServe(cfg.Address, nil)
}

func currentDir() (dir string) {
	dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	return
}
