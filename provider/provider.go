package provider

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"net/http"
	"bytes"
)

type Generic interface {
	GetAsset(string, string) ([]byte, error)
	WriteAsset(string, string, http.ResponseWriter) (error)
}

type Amazon struct {
	auth aws.Auth
	connection *s3.S3
}

func AWSConnect(access_key string, secret_key string) (amazon Amazon, err error) {
	var auth aws.Auth
	if auth, err = aws.GetAuth(access_key, secret_key); err != nil { return }

	connection := s3.New(auth, aws.USEast)
	amazon = Amazon{auth: auth, connection: connection}
	return
}

func (amazon Amazon) GetAsset(bucketName string, name string) ([]byte, error) {
	bucket := amazon.connection.Bucket(bucketName)
	return bucket.Get(name)
}

func (amazon Amazon) WriteAsset(bucketName string, assetName string, w http.ResponseWriter) (err error) {
	var data []byte
	if data, err = amazon.GetAsset(bucketName, assetName); err != nil { return }

	body := bytes.NewBuffer(data)
	buffer := make([]byte, 1024)

	for n, e := body.Read(buffer) ; e == nil ; n, e	= body.Read(buffer) {
		if n > 0 {
			w.Write(buffer[0:n])
		}
	}

	return
}
