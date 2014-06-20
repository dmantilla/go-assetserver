package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	AWS interface {}
	Address string
	Cache interface {}
}

func FileName() (filename string) {
	if fn := os.Getenv("FILES_ENV"); fn == "" {
		filename = "development.json"
	} else { filename = fn + ".json" }
	return
}

func Load(path string) (config Configuration, err error) {
	file, err := os.Open(path + "/" + FileName())
	if err != nil { return }

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	return
}

func (config *Configuration) AwsNode(key string) string {
	return config.AWS.(map[string]interface {})[key].(string)
}

func (config *Configuration) CacheNode(key string) string {
	return config.Cache.(map[string]interface {})[key].(string)
}
