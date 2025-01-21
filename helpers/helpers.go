package helpers

import (
	"log"
	"os"
	"strings"
)

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func CountFilesInDir(folder string) int {
	files, err := os.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() && !strings.HasSuffix(file.Name(), "exe") {
			count++
		}
	}
	return count
}

func Map[T any, U any](input []T, fn func(T) U) []U {
	output := make([]U, len(input))
	for i, v := range input {
		output[i] = fn(v)
	}
	return output
}

func KeyMap[T any, K comparable](input []T, fn func(T) K) map[K]T {
	m := make(map[K]T)
	for _, item := range input {
		m[fn(item)] = item
	}
	return m
}
