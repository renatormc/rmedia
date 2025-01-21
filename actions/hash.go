package actions

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"rmedia/helpers"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type HashResult struct {
	RelPath string
	Hash    string
	Err     error
}

func HashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hasher := sha512.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash), nil
}

func HashWorker(folder string, results chan<- HashResult, tasks <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		res, err := HashFile(task)
		relPath := helpers.Must(filepath.Rel(folder, task))
		results <- HashResult{Hash: res, Err: err, RelPath: relPath}
	}
}

func CountFiles(folder string) int64 {
	var total int64 = 0
	IterateFilesRecursively(folder, func(path string) {
		total += 1
	})
	return total
}

func IterateFilesRecursively(root string, handler func(string)) {
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if !d.IsDir() {
			handler(path)
		}
		return nil
	})
}

func HashProgressor(folder string, total int64, results <-chan HashResult, wg *sync.WaitGroup) {

	file, err := os.Create(filepath.Join(folder, "hash.txt"))
	if err != nil {
		panic(err)
	}
	defer func() {
		file.Close()
		wg.Done()
	}()

	writer := bufio.NewWriter(file)

	bar := progressbar.Default(total)

	for res := range results {
		if res.Err == nil {
			_, err := writer.WriteString(fmt.Sprintf("%s %s\n", res.Hash, strings.ReplaceAll(res.RelPath, "\\", "/")))
			if err != nil {
				panic(err)
			}
		} else {
			log.Printf("Error calculating hash from %q", res.RelPath)
		}

		bar.Add(1)
	}
	err = writer.Flush()
	if err != nil {
		panic(err)
	}

}
