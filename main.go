package main

import "fmt"

const (
	// MinAge is the minimum amount of days before a file gets deleted
	MinAge int = 30

	// MaxAge is the maximum amount of days before a file gets deleted
	MaxAge int = 365

	// MaxSize it the maximum size of a file (in bytes)
	MaxSize int64 = 512 * MiB

	// MiB is the size of 1 Mebibyte
	MiB = 1 << 20
)

func main() {
	fmt.Println("vim-go")
}
