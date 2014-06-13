package server

import (
	"net/http"
	"fmt"
	"os"
	"path/filepath"
	"github.com/gographics/imagick/imagick"
	"../config"
	"../provider"
)

var cfg config.Configuration
var amazon provider.Amazon

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	fmt.Println(AssetName(r))
	amazon.WriteAsset(cfg.AwsNode("assets_bucket"), AssetName(r), w, r)
}

func Run() {
	var err error

	imagick.Initialize()
	defer imagick.Terminate()

	dir := currentDir()
	cfg, err = config.Load(dir)
	if err != nil {
		fmt.Printf("There was a problem loading %s in %s: %s\n", config.FileName(), dir, err)
		os.Exit(1)
	}

	amazon, err = provider.AWSConnect(cfg.AwsNode("access_key"), cfg.AwsNode("secret_key"))
	if err != nil {
		fmt.Println("Cannot connect to Amazon")
		os.Exit(1)
	}

	http.HandleFunc("/", handler)
	fmt.Println("Running...")
	http.ListenAndServe(cfg.Address, nil)
}

func AssetName(r *http.Request) string {
	return r.URL.Path[1:]
}

func currentDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}
