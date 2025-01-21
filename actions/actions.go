package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"rmedia/helpers"
	"sync"
)

func Zip(folder string, mediaType string, compressLevel int) {
	folder = helpers.Must(filepath.Abs(folder))

	outFolder := filepath.Join(filepath.Dir(folder), "burn_medias")
	if helpers.DirectoryExists(outFolder) {
		helpers.CheckError(os.RemoveAll(outFolder))
	}
	helpers.CheckError(os.Mkdir(outFolder, os.ModePerm))
	if err := Compress7z(folder, outFolder, mediaType, compressLevel); err != nil {
		panic(err)
	}

	OrganizeFolders(outFolder)

}

func Hash(path string, workers int) {
	if helpers.FileExists(path) {
		res := helpers.Must(HashFile(path))
		fmt.Println(res)
	} else {
		fmt.Println("Counting files...")
		total := CountFiles(path)
		tasks := make(chan string)
		results := make(chan HashResult)
		var wg sync.WaitGroup
		var progressWg sync.WaitGroup

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go HashWorker(path, results, tasks, &wg)
		}
		progressWg.Add(1)
		go HashProgressor(path, total, results, &progressWg)

		IterateFilesRecursively(path, func(path string) {
			tasks <- path
		})
		close(tasks)
		wg.Wait()
		close(results)
		progressWg.Wait()
		hashHash := helpers.Must(HashFile(filepath.Join(path, "hash.txt")))
		fmt.Printf("Hash do hash:\n%s\n", hashHash)
	}

}
