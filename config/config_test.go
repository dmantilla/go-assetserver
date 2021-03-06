package config

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"os"
)

func TestFileName(t *testing.T) {
	Convey("it assumes 'development' by default", t, func() {
		os.Setenv("FILES_ENV", "")
		So(FileName(), ShouldEqual, "development.json")
	})

	Convey("it uses FILES_ENV as the filename", t, func() {
		os.Setenv("FILES_ENV", "test")
		So(FileName(), ShouldEqual, "test.json")
	})
}

func TestLoad(t *testing.T) {
	Convey("it fails loading the file", t, func() {
		os.Setenv("FILES_ENV", "unknown")
		_, err := Load("..")
		So(err, ShouldNotEqual, nil )
	})

	Convey("it loads the file", t, func() {
		os.Setenv("FILES_ENV", "configuration_sample")
		config, err := Load("..")
		aws := config.AWS.(map[string]interface{})

		So(err, ShouldEqual, nil )
		So(aws["access_key"], ShouldEqual, "access_key_value")
		So(aws["secret_key"], ShouldEqual, "secret_key_value")
		So(aws["assets_bucket"], ShouldEqual, "bucket_name")

		So(config.Address, ShouldEqual, ":4000")
	})
}

func TestAwsNode(t *testing.T) {
	Convey("return the right values", t, func() {
		os.Setenv("FILES_ENV", "configuration_sample")
		config, _ := Load("..")

		So(config.AwsNode("access_key"), ShouldEqual, "access_key_value")
		So(config.AwsNode("secret_key"), ShouldEqual, "secret_key_value")
		So(config.AwsNode("assets_bucket"), ShouldEqual, "bucket_name")
		So(config.AwsNode("rails_s3_url"), ShouldEqual, "http://s3.amazonaws.com/rails_bucket")
		So(config.AwsNode("legacy_s3_url"), ShouldEqual, "http://s3.amazonaws.com/legacy_bucket")
	})
}

func TestCacheNode(t *testing.T) {
	Convey("return the right values", t, func() {
			os.Setenv("FILES_ENV", "configuration_sample")
			config, _ := Load("..")

			So(config.CacheNode("folder"), ShouldEqual, "/tmp")
		})
}
