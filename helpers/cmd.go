package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func CmdExec(command string, args ...string) (*bytes.Buffer, error) {
	cmd := exec.Command(command, args...)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	cmd.Stderr = cmdOutput
	err := cmd.Run()
	if err != nil {
		return cmdOutput, err
	}
	return cmdOutput, nil
}

func CmdExecTextOutput(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	cmd.Stderr = cmdOutput
	cmd.Run()
	return strings.TrimSpace(cmdOutput.String())
}

func CmdExecConsole(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CmdExecStrOutput(command string, args ...string) (string, error) {
	res, err := CmdExec(command, args...)
	return res.String(), err
}

func Decode(codepage encoding.Encoding, data []byte) (string, error) {
	r := transform.NewReader(bytes.NewReader(data), codepage.NewDecoder())
	outUtf8, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(outUtf8), nil
}

func CmdExecStrOutputDecode(codepage encoding.Encoding, command string, args ...string) (string, error) {
	res, err := CmdExec(command, args...)
	if err != nil {
		return res.String(), err
	}
	b := res.Bytes()
	res2, err := Decode(codepage, b)
	if err != nil {
		return res.String(), err
	}
	return res2, err
}

func CmdExecJson(out any, command string, args ...string) error {
	data, err := CmdExec(command, args...)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data.Bytes(), out)
	if err != nil {
		return err
	}
	return nil
}

func PowershellExec(args ...string) error {
	ps, err := exec.LookPath("powershell.exe")
	if err != nil {
		return err
	}
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(ps, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
