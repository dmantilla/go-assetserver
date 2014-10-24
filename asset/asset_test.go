package asset

import (
	"testing"
	"net/http"
	"net/url"
	. "github.com/smartystreets/goconvey/convey"
)

func TestComputedName(t *testing.T) {
	Convey("valid scenarios", t, func() {
		a := New("/original/one.jpg", url.Values{}, nil, nil, nil)
		So(a.ComputedName(), ShouldEqual, "/original/one.jpg")
	})
	Convey("valid scenarios", t, func() {
		a := New("/original/one.jpg", url.Values{"w": []string{"25"}, "h": []string{"30"}}, nil, nil, nil)
		So(a.ComputedName(), ShouldEqual, "/original/one_25w_30h.jpg")
	})
}

func TestToBeResized(t *testing.T) {
	Convey("true", t, func() {
		q := url.Values{"w": []string{"21"}, "h": []string{"21"}}
		a := New("/original/one.jpg", q, nil, nil, nil)
		So(a.ToBeResized(), ShouldBeTrue)
	})

	Convey("false", t, func() {
		Convey("width is blank", func() {
			q := url.Values{"w": []string{""}, "h": []string{"30"}}
			a := New("/original/one.jpg", q, nil, nil, nil)
			So(a.ToBeResized(), ShouldBeTrue)
		})
		Convey("width is missing", func() {
			q := url.Values{"h": []string{"23"}}
			a := New("/original/one.jpg", q, nil, nil, nil)
			So(a.ToBeResized(), ShouldBeTrue)
		})
		Convey("height is blank", func() {
			q := url.Values{"w": []string{"40"}, "h": []string{""}}
			a := New("/original/one.jpg", q, nil, nil, nil)
			So(a.ToBeResized(), ShouldBeTrue)
		})
		Convey("height is missing", func() {
			q := url.Values{"w": []string{"25"}}
			a := New("/original/one.jpg", q, nil, nil, nil)
			So(a.ToBeResized(), ShouldBeTrue)
		})
		Convey("both are blank", func() {
			q := url.Values{"w": []string{""}, "h": []string{""}}
			a := New("/original/one.jpg", q, nil, nil, nil)
			So(a.ToBeResized(), ShouldBeFalse)
		})
		Convey("both are missing", func() {
			q := url.Values{}
			a := New("/original/one.jpg", q, nil, nil, nil)
			So(a.ToBeResized(), ShouldBeFalse)
		})
		Convey("dimensions are invalid", func() {
			q := url.Values{"w": []string{"20"}, "h": []string{"20"}}
			a := New("/original/one.jpg", q, nil, nil, nil)
			So(a.ToBeResized(), ShouldBeFalse)
		})
	})
}

func TestRequestedDimensions(t *testing.T) {
	Convey("Determine if an image has to be resized or not", t, func() {
		Convey("Returns zeroes when no dimensions are provided", func() {
			q := url.Values{}
			a := New("/original/one.jpg", q, nil, nil, nil)
			w, h := a.RequestedDimensions()
			So(w, ShouldEqual, 0)
			So(h, ShouldEqual, 0)
			})

		Convey("It accepts 'w' and 'h'", func() {
			q := url.Values{"w": []string{"10"}, "h": []string{"20"}}
			a := New("/original/one.jpg", q, nil, nil, nil)
			w, h := a.RequestedDimensions()
			So(w, ShouldEqual, 10)
			So(h, ShouldEqual, 20)
		})

		Convey("Using aliases for 'w' and 'h'", func() {
			q := url.Values{"width": []string{"100"}, "height": []string{"200"}}
			a := New("/original/one.jpg", q, nil, nil, nil)
			w, h := a.RequestedDimensions()
			So(w, ShouldEqual, 100)
			So(h, ShouldEqual, 200)
		})

	})
}

func TestSanitizeQueryParam(t *testing.T) {
	Convey("Converts height to proper values", t, func() {
			Convey("no value provided", func() {
					v, _ := SanitizeQueryParam("h", "")
					So(v, ShouldEqual, "0")
				})
			Convey("invalid value provided", func() {
					v, _ := SanitizeQueryParam("h", "_fiejf")
					So(v, ShouldEqual, "0")
				})
			Convey("another invalid value provided", func() {
					v, _ := SanitizeQueryParam("h", "a25")
					So(v, ShouldEqual, "0")
				})
			Convey("valid value provided", func() {
					v, _ := SanitizeQueryParam("h", "250")
					So(v, ShouldEqual, "250")
				})
		})
	Convey("Converts width to proper values", t, func() {
			Convey("no value provided", func() {
					v, _ := SanitizeQueryParam("w", "")
					So(v, ShouldEqual, "0")
				})
			Convey("invalid value provided", func() {
					v, _ := SanitizeQueryParam("w", "_fiejf")
					So(v, ShouldEqual, "0")
				})
			Convey("another invalid value provided", func() {
					v, _ := SanitizeQueryParam("w", "a25")
					So(v, ShouldEqual, "0")
				})
			Convey("valid value provided", func() {
					v, _ := SanitizeQueryParam("w", "250")
					So(v, ShouldEqual, "250")
				})
			Convey("width value provided", func() {
					v, _ := SanitizeQueryParam("width", "250")
					So(v, ShouldEqual, "250")
				})
			Convey("height value provided", func() {
					v, _ := SanitizeQueryParam("height", "125")
					So(v, ShouldEqual, "125")
				})
		})
	Convey("Returns error", t, func() {
			Convey("invalid parameter", func() {
					_, err := SanitizeQueryParam("unknown", "")
					So(err.Error(), ShouldEqual, "Parameter 'unknown' is not valid")
				})
		})
}

func TestQuery(t *testing.T) {
	Convey("Build a map with the valid keys", t, func() {
			Convey("No query parameters", func() {
					u := url.URL{RawQuery:""}
					r := http.Request{URL: &u }
					q := Query(&r)
					So(q.Get("h"), ShouldEqual, "")
					So(q.Get("w"), ShouldEqual, "")
				})

			Convey("Repeated and extra parameters", func() {
					u := url.URL{RawQuery:"w=50&m=30&h=35&unkown&h=80"}
					r := http.Request{URL: &u }
					q := Query(&r)
					So(q.Get("h"), ShouldEqual, "35")
					So(q.Get("w"), ShouldEqual, "50")
					So(q.Get("m"), ShouldEqual, "")
					So(q.Get("unknown"), ShouldEqual, "")
				})
		})
}
