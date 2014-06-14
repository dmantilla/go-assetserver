package provider

import (
	"testing"
	"net/http"
	"net/url"
	. "github.com/smartystreets/goconvey/convey"
)

func TestResizeRequired(t *testing.T) {
	Convey("Determine if an image has to be resized or not", t, func() {
		Convey("No query parameters, resize not required", func() {
			url := url.URL{RawQuery:""}
			r := http.Request{URL: &url }
			isRequired, _, _ := ResizeRequired(&r)
			So(isRequired, ShouldEqual, false)
		})

		Convey("Only one width, resize not required", func() {
			url := url.URL{RawQuery:"w=50"}
			r := http.Request{URL: &url }
			isRequired, _, _ := ResizeRequired(&r)
			So(isRequired, ShouldEqual, false)
		})

		Convey("Only one height, resize not required", func() {
			url := url.URL{RawQuery:"h=50"}
			r := http.Request{URL: &url }
			isRequired, _, _ := ResizeRequired(&r)
			So(isRequired, ShouldEqual, false)
		})

		Convey("Both dimensions, but one less than the minimum, don't resize", func() {
			url := url.URL{RawQuery:"w=100&h=20"}
			r := http.Request{URL: &url }
			isRequired, _, _ := ResizeRequired(&r)
			So(isRequired, ShouldEqual, false)
		})

		Convey("Both dimensions, resize!", func() {
			url := url.URL{RawQuery:"w=50&h=80"}
			r := http.Request{URL: &url }
			isRequired, w, h := ResizeRequired(&r)
			So(isRequired, ShouldEqual, true)
			So(w, ShouldEqual, 50)
			So(h, ShouldEqual, 80)
		})
	})
}
