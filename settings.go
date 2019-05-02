package main

import (
	"b2backup/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kurin/blazer/b2"
)

type PathDetails struct {
	Client         *b2.Client `json:"-"`
	B2Bucket       *b2.Bucket `json:"-"`
	AccountID      string
	ApplicationKey string
	Bucket         string
	Source         string
	Destination    string
	Compression    bool
	BackupSettings int16
	Signatures     []FileDetails
}

type FileDetails struct {
	Hash       string // Sha1 Hash of first 100MB of a file + path
	Path       string // Path inside Source path, excluding source path
	B2Path     string // Path inside bucket of backup
	Compressed bool   // Is the file compressed inside of b2
}

var pathList []string
var paths []PathDetails

// Load up .b2backup.json in path
func loadSettings() {

	var pathsChecked []string

	for _, path := range pathList {
		dat, errRead := ioutil.ReadFile(path + "/.b2backup.json")
		if errRead != nil {
			fmt.Println("There was a problem loading a path and it has been removed from the list: " + path)
			continue
		}

		// Keep it in the new checked list
		pathsChecked = append(pathsChecked, path)

		// Read the json
		var pathDetails PathDetails
		errUn := json.Unmarshal(dat, &pathDetails)
		if errUn != nil {
			fmt.Println("There was a problem loading a path json and it has been removed from the list: " + path)
			continue
		}

		// Load into memory
		paths = append(paths, pathDetails)
	}

	// Restore list to memory and write file
	pathList = pathsChecked
	savePaths()
}

func saveSettings(path PathDetails) {

	data, errMarsh := json.Marshal(path)
	if errMarsh != nil {
		log.Fatal(errMarsh)
	}

	errPaths := ioutil.WriteFile(path.Source+"/.b2backup.json", data, 0644)
	if errPaths != nil {
		log.Fatal(errPaths)
	}
}

// Load up paths.json in working directory
func loadPaths() error {

	wd, errWd := os.Getwd()
	if errWd != nil {
		return errWd
	}

	// No paths exists
	if utils.FileExists(wd+"/"+pathListFile) == false {
		return errors.New("No paths located in the system, try b2backup.run -init=<path>")
	}

	paths, errPaths := ioutil.ReadFile(wd + "/" + pathListFile)
	if errPaths != nil {
		return errPaths
	}

	if err := json.Unmarshal(paths, &pathList); err != nil {
		return errors.New("There was a problem parsing " + pathListFile)
	}

	// TODO: there are no paths inside path.json exit
	return nil
}

func savePaths() {

	wd, errWd := os.Getwd()
	if errWd != nil {
		log.Fatal(errWd)
	}

	data, errMarsh := json.Marshal(pathList)
	if errMarsh != nil {
		log.Fatal(errMarsh)
	}

	errPaths := ioutil.WriteFile(wd+"/"+pathListFile, data, 0644)
	if errPaths != nil {
		log.Fatal(errPaths)
	}
}

// Check if it exists already
func pathExists(path string) bool {

	for _, pathItem := range pathList {
		if pathItem == path {
			return true
		}
	}

	return false
}
