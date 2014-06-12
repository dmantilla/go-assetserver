package server

import (
	"net/http"
	"fmt"
	"os"
	"path/filepath"
	"bytes"
	"../config"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

var cfg config.Configuration
var err error
var awsAuth aws.Auth
var s3Connection *s3.S3
var s3Bucket *s3.Bucket

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(AssetName(r))
	data, ex := GetAsset(AssetName(r))
	if ex != nil { fmt.Fprint(w, "Resource not found") }

	body := bytes.NewBuffer(data)
	buffer := make([]byte, 1024)

	for n, e := body.Read(buffer) ; e == nil ; n, e	= body.Read(buffer) {
		if n > 0 {
			w.Write(buffer[0:n])
		}
	}
}

func Run() {
	dir := currentDir()
	cfg, err = config.Load(dir)
	if err != nil {
		fmt.Printf("There was a problem loading %s in %s: %s\n", config.FileName(), dir, err)
		os.Exit(1)
	}

	awsAuth, err = aws.GetAuth(cfg.AwsNode("access_key"), cfg.AwsNode("secret_key"))
	s3Connection = s3.New(awsAuth, aws.USEast)
	s3Bucket = s3Connection.Bucket(cfg.AwsNode("assets_bucket"))

	http.HandleFunc("/", handler)
	fmt.Println("Running...")
	http.ListenAndServe(":4000", nil)
}

func AssetName(r *http.Request) string {
	return r.URL.Path[1:]
}

func GetAsset(path string) (data []byte, e error) {
	return s3Bucket.Get(path)
}

func currentDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}
