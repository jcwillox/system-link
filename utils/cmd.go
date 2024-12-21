package utils

import (
	"bytes"
	"fmt"
	"github.com/google/shlex"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/execabs"
	"html/template"
	"runtime"
	"syscall"
)

type CommandConfig struct {
	Command string `json:"command"`
	Hidden  *bool  `json:"hidden"`
	Shell   string `json:"shell"`
}

func GetDefaultShell() []string {
	if runtime.GOOS == "windows" {
		return []string{"cmd", "/c"}
	}
	for _, shell := range []string{"bash", "ash", "sh"} {
		path, _ := execabs.LookPath(shell)
		if path != "" {
			return []string{path, "-c"}
		}
	}
	return []string{"sh", "-c"}
}

func getShell(shell string) []string {
	switch shell {
	case "cmd":
		return []string{"cmd", "/c"}
	case "powershell":
		return []string{"powershell", "-NoProfile", "-NoLogo", "-Command"}
	case "pwsh":
		return []string{"pwsh", "-NoProfile", "-NoLogo", "-Command"}
	case "":
		return GetDefaultShell()
	default:
		return []string{shell, "-c"}
	}
}

func renderTemplate(command string) (string, error) {
	exePath, exeDir, exeName, err := ExecutablePaths()
	if err != nil {
		return "", err
	}
	parse, err := template.New("command").Parse(command)
	if err != nil {
		return "", fmt.Errorf("failed to parse template command: %s; %w", command, err)
	}
	var tpl bytes.Buffer
	err = parse.Execute(&tpl, map[string]interface{}{
		"exePath": exePath,
		"exeDir":  exeDir,
		"exeName": exeName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to render template command: %s; %w", command, err)
	}
	return tpl.String(), nil
}

func RunCommand(command string, shell string, hidden *bool) error {
	var args []string
	var err error

	log.Info().Str("command", command).Str("shell", shell).
		Interface("hidden", hidden).Msg("running command")

	command, err = renderTemplate(command)
	if err != nil {
		return err
	}

	log.Info().Str("command", command).Msg("templated command")

	if shell == "none" {
		cmdArgs, err := shlex.Split(command)
		if err != nil {
			return fmt.Errorf("failed to parse command: %s; %w", command, err)
		}
		args = cmdArgs
	} else {
		shellCmd := getShell(shell)
		args = append(shellCmd, command)
	}

	cmd := execabs.Command(args[0], args[1:]...)

	if runtime.GOOS == "windows" && (hidden == nil || *hidden) {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}

	return cmd.Run()
}
