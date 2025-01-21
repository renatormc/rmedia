package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"rmedia/config"
	"rmedia/helpers"
	"runtime"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/text/encoding/charmap"
)

func Burn(folder string, diskName string, speed int) {
	cf := config.GetConfig()
	if runtime.GOOS == "windows" {
		dev := ChooseRecorder()
		cmd := exec.Command(
			cf.ExeCDBurnXP,
			"--burn-data",
			fmt.Sprintf("-device:%d", dev),
			fmt.Sprintf("-folder[\\]:%s", folder),
			fmt.Sprintf("-name:%s", diskName),
			"-verify",
			fmt.Sprintf("-speed:%d", speed),
			"-udf:1.02",
			"-close",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
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

type Recorder struct {
	Label  string
	Number int
}

func ListDevices() []Recorder {
	cf := config.GetConfig()
	output, err := helpers.CmdExecStrOutputDecode(charmap.CodePage850, cf.ExeCDBurnXP, "--list-drives")
	if err != nil {
		panic(err)
	}
	output = strings.TrimSpace(output)
	lines := strings.Split(output, "\n")
	return helpers.Map(lines, func(line string) Recorder {
		parts := strings.Split(line, ":")
		return Recorder{Label: line, Number: helpers.Must(strconv.Atoi(parts[0]))}
	})
}

func ChooseRecorder() int {
	recs := ListDevices()
	lines := helpers.Map(recs, func(rec Recorder) string { return rec.Label })
	km := helpers.KeyMap(recs, func(rec Recorder) string { return rec.Label })
	var device string
	prompt := &survey.Select{
		Message: "Selecione a gravadora:",
		Options: lines,
		Default: lines[0],
	}
	if err := survey.AskOne(prompt, &device); err != nil {
		panic(err)
	}
	return km[device].Number
}
