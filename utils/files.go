package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

// FileExists checks whether a file exists
func FileExists(f string) bool {
	_, err := os.Stat(f)
	if err != nil {
		return false
	}
	return true
}

// FileSha1 gets the sha1 signature
func FileSha1(path string) (string, error) {

	// Open the file
	file, errOpen := os.Open(path)
	if errOpen != nil {
		return "", errOpen
	}
	defer file.Close()

	// Prepare return and hash interface
	hash := sha1.New()

	// Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// Get the 20 bytes hash
	hashInBytes := hash.Sum(nil)[:20]

	// Return the bytes as a string
	return hex.EncodeToString(hashInBytes), nil
}
