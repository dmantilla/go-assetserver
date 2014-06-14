package provider

import (
	"os"
	"io/ioutil"
	"path/filepath"
	_ "fmt"
)

type CacheProvider struct {
	folderFullPath string
	entries []os.FileInfo
}

func CacheConnect(cacheFolder string) CacheProvider {
	return CacheProvider{folderFullPath: cacheFolder}
}

func (cache CacheProvider) GetFile(fileName string) (file []byte, err error) {
	fullName := filepath.Join(cache.folderFullPath, fileName)
	if _, err = os.Stat(fullName); err != nil { return }
	file, err = ioutil.ReadFile(fullName)
	return
}

func (cache CacheProvider) WriteFile(fileName string, data []byte) (error) {
	return ioutil.WriteFile(fileName, data, os.ModePerm)
}
