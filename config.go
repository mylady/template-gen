package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//FileReplceConfig : model for single file config
type FileReplceConfig struct {
	Template string `json:"template"`
	Replacer string `json:"replacer"`
	FileName string `json:"file_name"`
}

//AppConfig : model for app config
type AppConfig struct {
	Extension       string             `json:"extension"`
	PlaceHolder     string             `json:"place_holder"`
	DefaultTemplate string             `json:"default_template"`
	Files           []FileReplceConfig `json:"files"`
}

//ParseConfig : parse config from json
func ParseConfig(file string) (config *AppConfig, err error) {
	if _, err = os.Stat(file); err != nil {
		return nil, err
	}

	var content []byte
	if content, err = ioutil.ReadFile(file); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, err
}
