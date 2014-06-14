package provider

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"net/http"
	"bytes"
	"strconv"
	"../transformer"
)

type Amazon struct {
	auth aws.Auth
	connection *s3.S3
	cache CacheProvider
}

func AWSConnect(access_key string, secret_key string, cacheFolder string) (amazon Amazon, err error) {
	var auth aws.Auth
	if auth, err = aws.GetAuth(access_key, secret_key); err != nil { return }

	connection := s3.New(auth, aws.USEast)
	amazon = Amazon{auth: auth, connection: connection, cache: CacheConnect(cacheFolder)}
	return
}

func (amazon Amazon) GetAsset(bucketName string, name string) ([]byte, error) {
	bucket := amazon.connection.Bucket(bucketName)
	return bucket.Get(name)
}

func (amazon Amazon) ReadAsset(bucketName string, assetName string) (data []byte, err error) {
	// Try to grab the file from the cache
	if data, err = amazon.cache.GetFile(assetName); err != nil {
		// If not found go to Amazon
		data, err = amazon.GetAsset(bucketName, assetName)
	}
	return
}

func (amazon Amazon) WriteAsset(bucketName string, assetName string, response http.ResponseWriter, request *http.Request) (err error) {
	var data, image []byte

	if data, err = amazon.ReadAsset(bucketName, assetName); err != nil { return }

	resizeRequired, width, height := ResizeRequired(request)
	if resizeRequired {
		if image, err = transformer.Resize(data, width, height); err != nil { return }
		WriteData(image, response)
	} else {
		WriteData(data, response)
	}
	return
}

func WriteData(data []byte, response http.ResponseWriter) {
	body := bytes.NewBuffer(data)
	buffer := make([]byte, 1024)
	for n, e := body.Read(buffer) ; e == nil ; n, e	= body.Read(buffer) {
		if n > 0 {
			response.Write(buffer[0:n])
		}
	}
}

func ResizeRequired(r *http.Request) (isRequired bool, width uint, height uint) {
	var err error
	var w, h uint64

	values := r.URL.Query()
	if w, err = strconv.ParseUint(values.Get("w"), 10, 0); err != nil { w = 0 }
	if h, err = strconv.ParseUint(values.Get("h"), 10, 0); err != nil { h = 0 }
	width = uint(w)
	height = uint(h)
	isRequired = width > 20 && height > 20
	return
}
