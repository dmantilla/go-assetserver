package provider

import (
	"os"
	"io/ioutil"
	"path"
	"path/filepath"
)

type CacheProvider struct {
	folderFullPath string
	entries []os.FileInfo
}

func CacheConnect(cacheFolder string) CacheProvider {
	return CacheProvider{folderFullPath: cacheFolder}
}

func (cache CacheProvider) GetFile(fileName string) (file []byte, err error) {
	fullName := cache.FullFileName(fileName)
	if _, err = os.Stat(fullName); err != nil { return }
	file, err = ioutil.ReadFile(fullName)
	return
}

func (cache CacheProvider) WriteFile(fileName string, data []byte) (error) {
	fullName := cache.FullFileName(fileName)
	cache.CreateDirectories(fullName)
	return ioutil.WriteFile(fullName, data, 0755)
}

func (cache CacheProvider) CreateDirectories(fullFileName string) {
	os.MkdirAll(path.Dir(fullFileName), 0755)
}

func (cache CacheProvider) FullFileName(fileName string) string {
	return filepath.Join(cache.folderFullPath, fileName)
}
