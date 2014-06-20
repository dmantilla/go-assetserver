package provider

import (
	"os"
	"../config"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConnect(t *testing.T) {
	os.Setenv("FILES_ENV", "configuration_sample")
	cfg, _ := config.Load("..")

	Convey("Creates a new cache", t, func() {
		Convey("the new instance stores the folder name", func() {
			c := CacheConnect(cfg)
			So(c.folderFullPath, ShouldEqual, "/tmp")
		})
	})
}
