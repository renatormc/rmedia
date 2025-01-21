package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"rmedia/config"
	"rmedia/helpers"
	"runtime"
)

func Burn(folder string, diskName string, speed int, recorder int) {
	cf := config.GetConfig()
	if runtime.GOOS == "windows" {
		cmd := exec.Command(
			cf.ExeCDBurnXP,
			"--burn-data",
			fmt.Sprintf("-device:%d", recorder),
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

func ListRecorders() {
	helpers.CmdExecConsole(config.GetConfig().ExeCDBurnXP, "--list-drives")

}
