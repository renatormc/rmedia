package config

import (
	"os"
	"path/filepath"
	"rmedia/helpers"
	"strings"

	"github.com/joho/godotenv"
)

var config *Config

type Config struct {
	AppDir      string
	Exe7z       string
	ExeCDBurnXP string
}

func GetAppDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	dir := filepath.Dir(ex)

	if strings.Contains(dir, "go-build") {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return wd
	}
	return filepath.Dir(ex)
}

func GetConfig() *Config {
	if config == nil {
		config = &Config{AppDir: GetAppDir()}
		path := filepath.Join(config.AppDir, ".env")
		if helpers.FileExists(path) {
			err := godotenv.Load(filepath.Join(config.AppDir, ".env"))
			if err != nil {
				panic(err)
			}
		}

		config.Exe7z = os.Getenv("EXE7z")
		if config.Exe7z == "" {
			config.Exe7z = "7z"
		}
		config.ExeCDBurnXP = os.Getenv("ExeCDBurnXP")
	}
	return config
}
