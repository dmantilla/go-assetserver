package provider

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAWSConnect(t *testing.T) {
	Convey("tries to connect", t, func() {
		Convey("it fails", func() {
			_, err := AWSConnect("invalid", "credentials")
			So(err, ShouldNotEqual, nil)
		})
	})
}
