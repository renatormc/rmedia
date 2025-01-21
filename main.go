package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"rmedia/actions"
	"rmedia/config"

	"github.com/akamensky/argparse"
)

const EX = `rmedia burn --recorder 0 --speed 12 --name 123_2025 .
rmedia zip --max-size dvd --compress-level 3 .`

func main() {
	parser := argparse.NewParser("RMedia", "App managing medias.")
	logFile := parser.String("l", "log", &argparse.Options{Help: "Log file path"})

	burnCmd := parser.NewCommand("burn", "Burn disk")
	folderBurn := burnCmd.StringPositional(&argparse.Options{Help: "Folder to be burned.", Default: "."})
	speed := burnCmd.Int("s", "speed", &argparse.Options{Help: "Record speed", Default: 12})
	recorder := burnCmd.Int("r", "recorder", &argparse.Options{Help: "Recorder number", Default: 0})
	diskName := burnCmd.String("n", "name", &argparse.Options{Help: "Disk name", Required: true})

	listRecordersCmd := parser.NewCommand("list-recorders", "List recorders")

	zipCmd := parser.NewCommand("zip", "Zip folder")
	folderZip := zipCmd.StringPositional(&argparse.Options{Help: "Folder to be compressed", Default: "."})
	maxSize := zipCmd.String("s", "max-size", &argparse.Options{Help: "Max size", Default: "dvd"})
	compressLevel := zipCmd.Int("m", "compress-level", &argparse.Options{Help: "Compress level", Default: 3})

	hashCmd := parser.NewCommand("hash", "Hash folder")
	folderHash := hashCmd.StringPositional(&argparse.Options{Help: "Folder to be hashed", Default: "."})
	nWorkers := hashCmd.Int("w", "workers", &argparse.Options{Help: "Number of workers", Default: 3})

	configCmd := parser.NewCommand("config", "Open config file")

	examplesCmd := parser.NewCommand("ex", "Show command examples")

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
		actions.Burn(*folderBurn, *diskName, *speed, *recorder)
	case zipCmd.Happened():
		actions.Zip(*folderZip, *maxSize, *compressLevel)
	case hashCmd.Happened():
		actions.Hash(*folderHash, *nWorkers)
	case configCmd.Happened():
		path := filepath.Join(config.GetConfig().AppDir, ".env")
		cmd := exec.Command("notepad.exe", path)
		if err := cmd.Start(); err != nil {
			panic(err)
		}
	case listRecordersCmd.Happened():
		actions.ListRecorders()
	case examplesCmd.Happened():
		fmt.Println(EX)
	}
}
