package main

import (
	"math/rand"
	"os"
	"time"
)

// UploadFile takes in an array of bytes and lifetime in seconds and stores it
// to the fs, returning its unique name and any errors
func UploadFile(file []byte, lifetime time.Duration) (string, error) {
	f, err := os.Create(generateFileName(10))

	if err != nil {
		return "", err
	}

	f.Write(file)
	f.Close()

	return "", nil
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
