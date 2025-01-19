package main

import (
	"fmt"
	"log"
	"os"
	"rmedia/actions"
	"runtime"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("RMedia", "App managind medias")
	logFile := parser.String("l", "log", &argparse.Options{Help: "Log file path"})

	burnCmd := parser.NewCommand("burn", "Burn disk")
	folderBurn := burnCmd.StringPositional(&argparse.Options{Help: "Folder to be burned", Default: "."})

	zipCmd := parser.NewCommand("zip", "Zip folder")
	folderZip := burnCmd.StringPositional(&argparse.Options{Help: "Folder to be compressed", Default: "."})

	hashCmd := parser.NewCommand("hash", "Hash folder")
	folderHash := hashCmd.StringPositional(&argparse.Options{Help: "Folder to be hashed", Default: "."})
	nWorkers := hashCmd.Int("w", "workers", &argparse.Options{Help: "Number of workers", Default: runtime.NumCPU()})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if *logFile != "" {
		logFile, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	switch {

	case burnCmd.Happened():
		actions.Burn(*folderBurn)
	case zipCmd.Happened():
		actions.Zip(*folderZip)
	case hashCmd.Happened():
		actions.Hash(*folderHash, *nWorkers)

	}
}
