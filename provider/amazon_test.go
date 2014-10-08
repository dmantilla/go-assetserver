package provider

import (
	"testing"
	"log"
	"os"
	"../config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAmazon(t *testing.T) {
	os.Setenv("FILES_ENV", "configuration_sample")
	cfg, _ := config.Load("..")
	logger := log.New(os.Stdout, "", log.Ldate)
	amazon, _ := AWSConnect(cfg, logger)

	Convey("SourceURL", t, func() {
		Convey("it's a rails asset", func() {
			So(amazon.SourceURL("/regular/3929.jpg"), ShouldEqual, "http://s3.amazonaws.com/rails_bucket")
		})
		Convey("it's another asset", func() {
			So(amazon.SourceURL("/akd9939329.jpg"), ShouldEqual, "http://s3.amazonaws.com/legacy_bucket")
		})

	})
}
