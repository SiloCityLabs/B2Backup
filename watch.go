package main

import (
	"b2backup/utils"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func watchDirectory(path *PathDetails) {
	//apt-get install inotify-tools
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for run {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					fmt.Println("Watcher not ok (1)")
					return
				}

				switch event.Op.String() {
				case "WRITE": //CREATE,CHMOD
					fmt.Println(event.Op.String() + " " + event.Name)
					go b2NewFile(path, event.Name)
				case "REMOVE", "RENAME":
					fmt.Println(event.Op.String() + " " + event.Name)
					go b2RemoveFile(path, event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("Watcher not ok (2)")
					return
				}

				fmt.Println("Error: ", err)
			}
		}
	}()

	//TODO: rescan every x mins
	list := listDirectories(path.Source)
	for _, folder := range list {
		err = watcher.Add(folder)
		if err != nil {
			fmt.Println("There was an error startin up a path watcher: " + err.Error())
		} else {
			fmt.Println("Folder added: " + folder)
		}
	}

	<-done
}

func listDirectories(path string) []string {
	var list []string
	list = append(list, path)

	//Recursively scan folders
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Error listing directory: " + err.Error())
	}
	for _, f := range files {

		// Skip Hidden
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		//Add its children
		if f.IsDir() {

			list2 := listDirectories(path + "/" + f.Name())
			for _, folder := range list2 {
				list = append(list, folder)
			}
		}
	}

	return list
}

//TODO: prevent overlapping uploads of the same file (ctrl+s trigger happy person), prolly some queue with a 5 second wait
func b2NewFile(path *PathDetails, file string) {

	_, fileName := filepath.Split(file)

	// Skip Hidden
	if strings.HasPrefix(fileName, ".") {
		return
	}

	sha1, errSha := utils.FileSha1(file)
	if errSha != nil {
		fmt.Println(errSha)
		return
	}

	folderpath := strings.Replace(file, path.Source, "", 1)

	// Check if the file exists with the same signature
	for i := range path.Signatures {
		if path.Signatures[i].Path == folderpath && path.Signatures[i].Hash == sha1 {
			// Nothing about this file changed, lets go
			return
		}
	}

	// Prep the details
	fileDetails := FileDetails{
		B2Path:     path.Destination + folderpath,
		Compressed: false,
		Hash:       sha1,
		Path:       folderpath,
	}

	// Begin the upload
	if err := path.uploadFile(path.B2Bucket, file, strings.Replace(fileDetails.B2Path, "/", "", 1)); err != nil {
		fmt.Println("Something went wrong uploading file(" + file + "): " + err.Error())
		return
	}

	// Lets add it
	path.Signatures = append(path.Signatures, fileDetails)
	saveSettings(*path)
}

func b2RemoveFile(path *PathDetails, file string) {
	//TODO: all of this
}
