package main

import (
	_ "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "log"
	"os"
	_ "time"
)

func main() {

	fmt.Println("Thanks for using the Auto-Backup tool! This tool will automatically backup your files to a remote server.")
	fmt.Println("Trying to read config file [config.json]...")

	jsonFile, err := os.Open("users.json")
	if err != nil {
		// handle error
		// if the file doesn't exist, create it
		fmt.Println("Config file not found. Creating new config file...")
		// create a new file
		newFile, err := os.Create("users.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		// write to the file with the Config struct with pretty print
		prettyJSON, err := json.MarshalIndent(Config{
			RemotePath: "",
			LocalPaths: []string{"/Users/max/Downloads", "/Users/max/Downloads"},
		}, "", "    ")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = newFile.Write(prettyJSON)
	}
	fmt.Println("Successfully opened users.json")
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	ZipWriter(".idea/", "test")
	for i := 0; i < len(config.LocalPaths); i++ {
		fmt.Println(config.LocalPaths[i])
	}

}

type Config struct {
	RemotePath string   `json:"remote_path"`
	LocalPaths []string `json:"local_path"`
}
