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

func Zip(folder string) {
	if runtime.GOOS == "windows" {
		fmt.Println("Not implemented for Windows yet")
	} else {
		fmt.Println("Not implemented for Linux yet")
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
