package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	AWS interface {}
	Address string
}

func FileName() string {
	if fn := os.Getenv("FILES_ENV"); fn == "" { return "development.json" } else { return fn + ".json" }
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
