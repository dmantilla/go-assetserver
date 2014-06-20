package provider

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"net/http"
	"log"
	"io/ioutil"
	"strings"
	"fmt"
	"../config"
)

type Amazon struct {
	auth aws.Auth
	connection *s3.S3
	logger *log.Logger
	railsS3URL string
	legacyS3URL string
}

func AWSConnect(cfg config.Configuration, logger *log.Logger) (amazon Amazon, err error) {
	var auth aws.Auth
	if auth, err = aws.GetAuth(cfg.AwsNode("access_key"), cfg.AwsNode("secret_key")); err != nil { return }

	connection := s3.New(auth, aws.USEast)
	amazon = Amazon{
		auth: auth,
		connection: connection,
		logger: logger,
		railsS3URL: cfg.AwsNode("rails_s3_url"),
		legacyS3URL: cfg.AwsNode("legacy_s3_url"),
	}
	return
}

func (amazon Amazon) GetAsset(bucketName string, name string) ([]byte, error) {
	bucket := amazon.connection.Bucket(bucketName)
	return bucket.Get(name)
}

func (amazon Amazon) SourceURL(assetName string) (sourceUrl string) {
	if strings.Index(assetName, "/original") == 0 {
		sourceUrl = amazon.railsS3URL
	} else {
		sourceUrl = amazon.legacyS3URL
	}
	return
}

func (amazon Amazon) FetchAsset(assetName string) (data []byte, err error) {
	source := amazon.SourceURL(assetName)
	var response *http.Response
	resource := source + assetName
	amazon.logger.Printf("Retrieving %s", resource)
	if response, err = http.Get(resource); err == nil {
		if response.StatusCode == 200 {
			defer response.Body.Close()
			data, err = ioutil.ReadAll(response.Body)
		} else {
			err = fmt.Errorf("Resource %s not found", resource)
		}
	}
	return
}
