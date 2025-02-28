package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"rmedia/config"
	"rmedia/helpers"
	"runtime"
	"strings"

	"github.com/dustin/go-humanize"
)

const (
	DVD_SIZE                 = 4508876800
	CD_SIZE                  = 702545920
	BLUE_RAY_SIZE            = 23622320128
	BLUE_RAY_DUAL_LAYER_SIZE = 47244640256
	PEN_DRIVE_16_SIZE        = 15032385536
	PEN_DRIVE_32_SIZE        = 30064771072
	PEN_DRIVE_64_SIZE        = 65498251264
	PEN_DRIVE_128_SIZE       = 118111600640
)

func getMaxSize(value string) uint64 {
	switch value {
	case "dvd":
		return DVD_SIZE
	case "cd":
		return CD_SIZE
	case "blue-ray":
		return BLUE_RAY_SIZE
	case "blue-ray-dual":
		return BLUE_RAY_DUAL_LAYER_SIZE
	case "pen-drive-16":
		return PEN_DRIVE_16_SIZE
	case "pen-drive-32":
		return PEN_DRIVE_32_SIZE
	case "pen-drive-64":
		return PEN_DRIVE_64_SIZE
	case "pen-drive-128":
		return PEN_DRIVE_128_SIZE
	}

	return helpers.Must(humanize.ParseBytes(value))
}

func Compress7z(folder string, outputFolder string, maxSize string, compressLevel int) error {
	fName := "dados.7z"
	args := []string{"a", fmt.Sprintf("-v%d", getMaxSize(maxSize)), fmt.Sprintf("-mx%d", compressLevel)}
	if runtime.GOOS == "windows" {
		args = append(args, "-sfx7z.sfx")
		fName = "dados.exe"
	}
	args = append(args, filepath.Join(outputFolder, fName), folder)
	cmd := exec.Command(config.GetConfig().Exe7z, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func TutorialLinux(nMedias int) string {
	if nMedias == 1 {
		return "<h4>Copie o arquivo dados.7z.001 para um local de sua preferência e o descompacte utilizando o programa 7zip.</h4>"
	}
	files := []string{}
	for i := 0; i < nMedias; i++ {
		files = append(files, fmt.Sprintf("<li>dados.7z.%03d</li>", i+1))
	}
	aux := strings.Join(files, "\n")
	text := fmt.Sprintf("<h4>Crie uma nova pasta em um local de sua preferência e copie todos os arquivos abaixo. Em seguida abra o arquivo dados.7z001 no programa 7zip e execute a descompactação.</h4><ul>%s</ul>", aux)
	return text
}

func TutorialWindows(nMedias int) string {
	if nMedias == 1 {
		return "<h4>De dois cliques no arquivo dados.exe e descompacte para um local de sua preferência.</h4>"
	}
	files := []string{"<li>dados.exe</li>"}
	for i := 0; i < nMedias; i++ {
		files = append(files, fmt.Sprintf("<li>dados.7z.%03d</li>", i+1))
	}
	aux := strings.Join(files, "\n")
	text := fmt.Sprintf("<h4>Crie uma nova pasta em um local de sua preferência e copie todos os arquivos abaixo, copiando de todas as mídias. Em seguida dê dos clique no arquivo dados.exe e execute a descompactação.</h4><ul>%s</ul>", aux)
	return text
}

func Tutorial(nMedias int) string {
	if runtime.GOOS == "windows" {
		return TutorialWindows(nMedias)
	}
	return TutorialLinux(nMedias)

}

func OrganizeFolders(outFolder string) {
	nFiles := helpers.CountFilesInDir(outFolder)
	html := Tutorial(nFiles)
	if nFiles == 1 {
		helpers.CheckError(os.WriteFile(filepath.Join(outFolder, "intruções.html"), []byte(html), os.ModePerm))
	} else {
		for _, file := range helpers.Must(os.ReadDir(outFolder)) {
			var mediaFolder string
			if strings.HasSuffix(file.Name(), "exe") {
				mediaFolder = filepath.Join(outFolder, "001")
			} else {
				mediaFolder = filepath.Join(outFolder, file.Name()[len(file.Name())-3:len(file.Name())])
			}
			if !helpers.DirectoryExists(mediaFolder) {
				helpers.CheckError(os.Mkdir(mediaFolder, os.ModePerm))
			}

			helpers.CheckError(os.Rename(filepath.Join(outFolder, file.Name()), filepath.Join(mediaFolder, file.Name())))
			helpers.CheckError(os.WriteFile(filepath.Join(mediaFolder, "intruções.html"), []byte(html), os.ModePerm))
		}
	}
}
