package provider

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConnect(t *testing.T) {
	Convey("Creates a new cache", t, func() {
		Convey("the new instance stores the folder name", func() {
			c := CacheConnect("/tmp")
			So(c.folderFullPath, ShouldEqual, "/tmp")
		})
	})
}
