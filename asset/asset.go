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
	"../transformer"
)

type Asset struct {
	path string
	query url.Values

	data   []byte

	amazon *provider.Amazon
	cache  *provider.CacheProvider
}

func WriteResponse(request *http.Request, response http.ResponseWriter, amazon *provider.Amazon, cache *provider.CacheProvider) (err error) {
	a := Build(request, amazon, cache)
	err = a.Fetch()
	a.Write(response)
	return
}

func Build(request *http.Request, amazon *provider.Amazon, cache *provider.CacheProvider) Asset {
	return New(request.URL.Path, Query(request), amazon, cache)
}

func New(path string, query url.Values, amazon *provider.Amazon, cache *provider.CacheProvider) Asset {
	return Asset{path: path, query: query, amazon: amazon, cache: cache}
}

func (a *Asset) Write(response http.ResponseWriter) {
	body := bytes.NewBuffer(a.data)
	log.Printf("Data size: %d\n", len(a.data))
	buffer := make([]byte, 1024)
	for n, e := body.Read(buffer) ; e == nil ; n, e	= body.Read(buffer) {
		if n > 0 {
			response.Write(buffer[0:n])
		}
	}
}
func (a *Asset) FetchOriginalFromCloud() (err error) {
	if a.data, err = a.amazon.FetchAsset(a.path); err == nil {
		err = a.cache.WriteFile(a.path, a.data)
	}
	return
}

func (a *Asset) FetchOriginal() (err error) {
	if a.data, err = a.cache.GetFile(a.path); err != nil {
		err = a.FetchOriginalFromCloud()
	}
	return
}

func (a *Asset) Fetch() (err error) {
	// Look for requested file in cache
	computedName := a.ComputedName()
	if a.data, err = a.cache.GetFile(computedName); err != nil {
		if a.ToBeResized() {
			err = a.FetchOriginal()
			if err == nil {
				w, h := a.RequestedDimensions()
				if a.data, err = transformer.Resize(a.data, w, h); err == nil {
					err = a.cache.WriteFile(computedName, a.data)
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
		result = fmt.Sprintf("%s%s_%sx%s%s", dir, name, a.query.Get("h"), a.query.Get("w"), ext)
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
	w, _ = strconv.ParseUint(a.query.Get("w"), 10, 0)
	h, _ = strconv.ParseUint(a.query.Get("h"), 10, 0)
	return uint(w), uint(h)
}

func SanitizeQueryParam(paramName string, value string) (result string, err error) {
	var unsigned_int uint64

	switch paramName {
	case "h", "w":
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

