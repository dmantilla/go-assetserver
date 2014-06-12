package server

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"net/http"
)

func TestAssetName(t *testing.T) {
	Convey("gets the name", t, func() {
		url := url.URL{Path:"/abcde"}
		req := http.Request{URL: &url}
		So(AssetName(&req), ShouldEqual, "abcde")

		url.Path = "/"
		So(AssetName(&req), ShouldEqual, "")
	})
}
