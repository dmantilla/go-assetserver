package asset

import (
	"net/http"
	"net/url"
	"strconv"
	"fmt"
	"path"
	"strings"
	"bytes"
	"log"
	"../provider"
	"github.com/gographics/imagick/imagick"
)

type Asset struct {
	path   string
	query  url.Values

	data   []byte

	amazon *provider.Amazon
	cache  *provider.CacheProvider
	logger *log.Logger
}

func WriteResponse(request *http.Request, response http.ResponseWriter, amazon *provider.Amazon, cache *provider.CacheProvider, logger *log.Logger) (err error) {
	a := Build(request, amazon, cache, logger)
	err = a.Fetch()
	a.Write(response)
	return
}

func Build(request *http.Request, amazon *provider.Amazon, cache *provider.CacheProvider, logger *log.Logger) Asset {
	return New(request.URL.Path, Query(request), amazon, cache, logger)
}

func New(path string, query url.Values, amazon *provider.Amazon, cache *provider.CacheProvider, logger *log.Logger) Asset {
	return Asset{path: path, query: query, amazon: amazon, cache: cache, logger: logger}
}

func (a *Asset) Write(response http.ResponseWriter) {
	body := bytes.NewBuffer(a.data)
	buffer := make([]byte, 1024)
	for n, e := body.Read(buffer) ; e == nil ; n, e	= body.Read(buffer) {
		if n > 0 {
			response.Write(buffer[0:n])
		}
	}
}
func (a *Asset) FetchOriginalFromCloud() (err error) {
	if a.data, err = a.amazon.FetchAsset(a.path); err == nil {
		a.logger.Printf("Fetched original %s from cloud", a.path)
		err = a.cache.WriteFile(a.path, a.data)
		if err == nil { a.logger.Printf("%s cached", a.path) }
	}
	return
}

func (a *Asset) FetchOriginal() (err error) {
	if a.data, err = a.cache.GetFile(a.path); err == nil {
		a.logger.Printf("Fetched original %s from cache", a.path)
	} else {
		err = a.FetchOriginalFromCloud()
	}
	return
}

func (a *Asset) Fetch() (err error) {
	// Look for requested file in cache
	computedName := a.ComputedName()
	if a.data, err = a.cache.GetFile(computedName); err == nil {
		a.logger.Printf("Serving %s from cache", computedName)
	} else {
		if a.ToBeResized() {
			err = a.FetchOriginal()
			if err == nil {
				w, h := a.RequestedDimensions()
				if a.data, err = a.Resize(w, h); err == nil {
					err = a.cache.WriteFile(computedName, a.data)
					a.logger.Printf("%s cached", computedName)
				}
			}
		} else {
			err = a.FetchOriginalFromCloud()
		}
	}
	return
}

func (a *Asset) ComputedName() (result string) {
	dir, base := path.Split(a.path)
	name := strings.Split(base, ".")[0]
	ext := path.Ext(a.path)
	if a.ToBeResized() {
		w, h := a.RequestedDimensions()
		result = fmt.Sprintf("%s%s_%dw_%dh%s", dir, name, w, h, ext)
	} else {
		result = fmt.Sprintf("%s%s%s", dir, name, ext)
	}
	return
}

func (a *Asset) ToBeResized() bool {
	w, h := a.RequestedDimensions()
	return w > 20 && h > 20
}

func (a *Asset) RequestedDimensions() (uint, uint) {
	var w, h uint64
	qW := a.query.Get("w")
	qH := a.query.Get("h")
	if len(qW) == 0 { qW = a.query.Get("width") }
	if len(qH) == 0 { qH = a.query.Get("height") }

	w, _ = strconv.ParseUint(qW, 10, 0)
	h, _ = strconv.ParseUint(qH, 10, 0)
	return uint(w), uint(h)
}

func (a *Asset) Resize(width uint, height uint) (image []byte, err error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err = mw.ReadImageBlob(a.data); err != nil { return }
	if err = mw.ResizeImage(width, height, imagick.FILTER_LANCZOS, 1); err != nil { return }
	if err = mw.SetImageCompressionQuality(95); err != nil { return }

	image = mw.GetImageBlob()
	a.logger.Printf("%s resized to %dx%d", a.path, height, width)
	return
}

func SanitizeQueryParam(paramName string, value string) (result string, err error) {
	var unsigned_int uint64

	switch paramName {
	case "h", "w", "height", "width":
		if unsigned_int, err = strconv.ParseUint(value, 10, 0); err == nil {
			result = strconv.FormatUint(unsigned_int, 10)
		} else { result = "0" }
	default:
		err = fmt.Errorf("Parameter '%s' is not valid", paramName)
	}
	return
}

func Query(request *http.Request) url.Values {
	query := url.Values{}
	rQuery := request.URL.Query()
	for k, _ := range rQuery {
		value, _ := SanitizeQueryParam(k, rQuery.Get(k))
		query.Set(k, value)
	}
	return query
}

