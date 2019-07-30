package main

import (
	"crypto/sha256"
	"math/rand"
	"os"
	"time"

	"go.etcd.io/bbolt"
)

// UploadFile takes in an array of bytes and lifetime in seconds and stores it
// to the fs, returning its unique name and any errors
func UploadFile(file []byte, lifetime time.Duration, extension string, bucket *bbolt.Bucket) (string, error) {
	hash := sha256.Sum256(file)
	name := generateFileName(10) + extension
	f, err := os.Create("files/" + name)

	if err != nil {
		return "", err
	}

	f.Write(file)
	f.Close()

	return name, nil
}

// https://medium.com/@kpbird/golang-generate-fixed-size-random-string-dd6dbd5e63c0
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateFileName(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]

	}

	return string(b)
}
