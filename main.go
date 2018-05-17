package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	argsCnt := len(os.Args)
	if argsCnt != 2 {
		printHelp()
		return
	}

	var config *AppConfig
	var err error

	if config, err = parseConfig(os.Args[1]); err != nil {
		printError(err)
		return
	}

	if config == nil {
		printError(errors.New("config is empty"))
		return
	}

	//process config
	processConfig(config)

	printLog("process finished")
}

//FileReplceConfig : model for single file config
type FileReplceConfig struct {
	Template string `json:"template"`
	Replacer string `json:"replacer"`
	FileName string `json:"file_name"`
}

//AppConfig : model for app config
type AppConfig struct {
	PlaceHolder     string             `json:"place_holder"`
	DefaultTemplate string             `json:"default_template"`
	Files           []FileReplceConfig `json:"files"`
}

func parseConfig(file string) (config *AppConfig, err error) {
	config = &AppConfig{}

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

func processConfig(config *AppConfig) {
	var err error
	var defaultTemplate string

	if config.DefaultTemplate != "" {
		if defaultTemplate, err = readTextFile(config.DefaultTemplate); err != nil {
			printLog("read default template error")
			printError(err)
			return
		}
	}

	if config.PlaceHolder == "" {
		printLog("place holder is empty")
		return
	}

	var targetTemplate string
	var targetContent string
	for _, file := range config.Files {
		if file.Template == "" && defaultTemplate != "" {
			targetTemplate = defaultTemplate
		} else {
			if targetTemplate, err = readTextFile(file.Template); err != nil {
				printLog("file " + file.FileName + " generate failed because of template error")
				printError(err)
				continue
			}
		}

		if file.Replacer == "" {
			printLog("file " + file.FileName + " replacer is empty")
			continue
		}

		if file.FileName == "" {
			printLog("file " + file.FileName + " output file name is empty")
			continue
		}

		targetContent = strings.Replace(targetTemplate, config.PlaceHolder, file.Replacer, -1)
		if err = writeTextFile(targetContent, file.FileName); err != nil {
			printLog("write to output file " + file.FileName + " failed")
			printError(err)
		} else {
			printLog("write to output file " + file.FileName + " succeed")
		}
	}
}

func readTextFile(file string) (data string, err error) {
	var buffer []byte
	if buffer, err = ioutil.ReadFile(file); err == nil {
		data = string(buffer)
	}
	return data, err
}

func writeTextFile(data string, file string) (err error) {
	buffer := ([]byte)(data)
	err = ioutil.WriteFile(file, buffer, 0644)
	return err
}

func printHelp() {
	fmt.Println("usage template-gen file.json")
	fmt.Println("file.json should be placed beside template-gen")
}

func printLog(log string) {
	fmt.Println(log)
}

func printError(err error) {
	fmt.Println("error message:", err.Error())
}
