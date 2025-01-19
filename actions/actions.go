package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"rmedia/helpers"
	"runtime"
	"sync"
)

func Burn(folder string) {
	if runtime.GOOS == "windows" {
		fmt.Println("Not implemented for Windows yet")
	} else {
		absPath := helpers.Must(filepath.Abs(folder))
		cmd := exec.Command("growisofs", "-Z", "/dev/sr0", "-R", "-J", absPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

}

func Zip(folder string, mediaType string, compressLevel int) {
	folder = helpers.Must(filepath.Abs(folder))

	outputFolder := filepath.Join(filepath.Dir(folder), "burn_medias")
	if helpers.DirectoryExists(outputFolder) {
		helpers.CheckError(os.RemoveAll(outputFolder))
	}
	helpers.CheckError(os.Mkdir(outputFolder, os.ModePerm))

	if runtime.GOOS == "windows" {
		panic("Not implemented for Windows yet")
	} else {
		if err := Compress7z(folder, outputFolder, mediaType, compressLevel); err != nil {
			panic(err)
		}
	}
	nFiles := helpers.CountFilesInDir(outputFolder)
	html := Tutorial(nFiles)
	if nFiles == 1 {
		helpers.CheckError(os.WriteFile(filepath.Join(outputFolder, "intruções.html"), []byte(html), os.ModePerm))
	} else {
		for _, file := range helpers.Must(os.ReadDir(outputFolder)) {
			mediaFolder := filepath.Join(outputFolder, file.Name()[len(file.Name())-3:len(file.Name())])
			helpers.CheckError(os.Mkdir(mediaFolder, os.ModePerm))
			helpers.CheckError(os.Rename(filepath.Join(outputFolder, file.Name()), filepath.Join(mediaFolder, file.Name())))
			helpers.CheckError(os.WriteFile(filepath.Join(mediaFolder, "intruções.html"), []byte(html), os.ModePerm))
		}
	}

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
	}

}
