package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kurin/blazer/b2"
)

//TODO: possibly a bug, it shows (2) when its uploaded again. Need to keep a single version.
func (path *PathDetails) uploadFile(bucket *b2.Bucket, src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	obj := bucket.Object(dst)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, f); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

func (path *PathDetails) printObjects(bucket *b2.Bucket) error {
	iterator := bucket.List(ctx)
	for iterator.Next() {
		fmt.Println(iterator.Object())
	}
	return iterator.Err()
}

func (path *PathDetails) connect() {
	var err error
	path.Client, err = b2.NewClient(ctx, path.AccountID, path.ApplicationKey)
	if err != nil {
		log.Fatal("Error connecting to B2 for path '" + path.Source + "': " + err.Error())
	}

	buckets, errB := path.Client.ListBuckets(ctx)
	if errB != nil {
		log.Fatal("Error connecting to B2 for path '" + path.Source + "': " + errB.Error())
	}

	// Check if bucket exists
	for _, bucket := range buckets {
		if bucket.Name() == path.Bucket {
			path.B2Bucket = bucket
			return
		}
	}

	log.Fatal("Bucket '" + path.Bucket + "' does not exist for path '" + path.Source + "'")
}
