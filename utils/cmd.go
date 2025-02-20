package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/shlex"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/execabs"
	"html/template"
	"os"
	"os/exec"
	"runtime"
)

type CommandConfig struct {
	Command    string            `json:"command" validate:"required"`
	Hidden     *bool             `json:"hidden"`
	Shell      string            `json:"shell"`
	ShowErrors bool              `json:"show_errors"`
	ShowOutput bool              `json:"show_output"`
	Detached   bool              `json:"detached"`
	Env        map[string]string `json:"env"`
}

type CommandResult struct {
	Stdout []byte
	Stderr []byte
	Code   int
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
	parse, err := template.New("command").Parse(command)
	if err != nil {
		return "", fmt.Errorf("failed to parse template command: %s; %w", command, err)
	}
	var tpl bytes.Buffer
	err = parse.Execute(&tpl, map[string]interface{}{
		"ExePath":      ExePath,
		"ExeDirectory": ExeDirectory,
		"ExeName":      ExeName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to render template command: %s; %w", command, err)
	}
	return tpl.String(), nil
}

// RunCommand runs a specified command,
// if Detached is true the command result will be empty
func RunCommand(cfg CommandConfig) (CommandResult, error) {
	var args []string

	log.Info().Str("command", cfg.Command).Str("shell", cfg.Shell).
		Interface("hidden", cfg.Hidden).Msg("running command")

	command, err := renderTemplate(cfg.Command)
	if err != nil {
		return CommandResult{}, err
	}

	log.Info().Str("command", command).Msg("templated command")

	if cfg.Shell == "none" {
		cmdArgs, err := shlex.Split(command)
		if err != nil {
			return CommandResult{}, fmt.Errorf("failed to parse command: %s; %w", command, err)
		}
		args = cmdArgs
	} else {
		shellCmd := getShell(cfg.Shell)
		args = append(shellCmd, command)
	}

	cmd := execabs.Command(args[0], args[1:]...)

	if len(cfg.Env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range cfg.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	if cfg.Hidden == nil || *cfg.Hidden {
		makeCmdHidden(cmd)
	}

	if cfg.Detached {
		if err := cmd.Start(); err != nil {
			return CommandResult{}, fmt.Errorf("failed to start command: %s; %w", command, err)
		}
		err := cmd.Process.Release()
		if err != nil {
			return CommandResult{}, err
		}
		return CommandResult{}, nil
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	res := CommandResult{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
		Code:   0,
	}

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		res.Code = ee.ExitCode()

		if cfg.ShowErrors {
			log.Err(err).Int("exit", ee.ExitCode()).Str("stdout", string(res.Stdout)).Str("stderr", string(res.Stderr)).Msg("command failed")
		}

		return res, nil
	}

	if cfg.ShowOutput {
		log.Info().Str("stdout", string(res.Stdout)).Str("stderr", string(res.Stderr)).Msg("command output")
	}

	return res, err
}
