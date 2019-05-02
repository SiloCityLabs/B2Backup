package main

import (
	"fmt"
	"log"
)

func newPath(path string) {
	var settings PathDetails

	fmt.Println("===============================")
	fmt.Println("== Starting new init")
	fmt.Println("===============================")

	// Load existing paths first
	if err := loadPaths(); err != nil {
		fmt.Println(err.Error())
		fmt.Println("There was a problem loading existing paths file, if you continue it will start fresh. Ctrl+C to cancel")
	}

	// Check if it exists and tell user
	if pathExists(path) {
		fmt.Println("This path already exists, if you continue it will start fresh. Ctrl+C to cancel")
	} else {
		pathList = append(pathList, path)
	}

	// Defaults
	settings.Source = path

	// Defaults that we are not supporting yet
	settings.Compression = false
	settings.BackupSettings = 1

	//TODO: Validation
	// Ask the user for details about the path
	settings.AccountID = askQuestion("Enter your B2 Account ID: ")
	settings.ApplicationKey = askQuestion("Enter your B2 Application Key: ")
	settings.Bucket = askQuestion("Enter the bucket you would like to backup to: ")
	settings.Destination = askQuestion("Enter the path you would like to backup to inside the bucket: ")

	// Check if we are able to connect
	settings.connect()

	// Save this settings
	saveSettings(settings)

	//Save paths file
	savePaths()

	fmt.Println("Path has been created, Start the program as a service now")
}

func askQuestion(question string) string {

	fmt.Print(question)

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	return response
}
