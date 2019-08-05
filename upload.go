package main

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"go.etcd.io/bbolt"
)

// UploadFile takes in an array of bytes and lifetime in seconds and stores it
// to the fs, returning its unique name and any errors
func UploadFile(file []byte, fileSize int64, extension string, bucket *bbolt.Bucket) (string, error) {
	hash := sha256.Sum256(file)
	val := bucket.Get(hash[:])

	// If there was an entry already, return it (file already exists)
	if val != nil {
		file, err := os.Stat("files/" + string(val))

		if err != nil {
			fmt.Println(err)
		}

		// TODO: Race condition? what if the file is being deleted as we are sending it back?
		// Maybe use a channel to block here?
		timePassed := float64(daysBetween(file.ModTime(), time.Now()))
		maxAge := math.Floor(calculateAge(MinAge, MaxAge, file.Size(), MaxSize))

		// If it's ~95% through its lifetime, reupload
		if !(timePassed/maxAge <= 0.95) {
			return string(val), nil
		}

		// Delete as it's going to get put again afterwards
		bucket.Delete(hash[:])
	}

	name := generateFileName(10) + extension

	f, err := os.Create("files/" + name)

	if err != nil {
		return "", err
	}

	f.Write(file)
	f.Close()

	bucket.Put(hash[:], []byte(name))

	return name, nil
}

// https://medium.com/@kpbird/golang-generate-fixed-size-random-string-dd6dbd5e63c0
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// generateFileName takes in a length and generates a random name
func generateFileName(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]

	}

	return string(b)
}
