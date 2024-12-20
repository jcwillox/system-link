package utils

import (
	"errors"
	"github.com/shirou/gopsutil/v4/process"
	"os"
	"path/filepath"
	"strconv"
)

type InstanceLock struct {
	lockFile string
}

func NewInstanceLock() (InstanceLock, error) {
	exePath, err := os.Executable()
	if err != nil {
		return InstanceLock{}, err
	}
	return InstanceLock{lockFile: filepath.Join(exePath + ".lock")}, nil
}

// Lock writes the current process id to the lock file
func (l InstanceLock) Lock() error {
	return os.WriteFile(l.lockFile, []byte(strconv.Itoa(os.Getpid())), 0644)
}

// Unlock removes the lock file
func (l InstanceLock) Unlock() error {
	return os.Remove(l.lockFile)
}

// LockedPid returns the pid in the lock file
func (l InstanceLock) LockedPid() (int, error) {
	data, err := os.ReadFile(l.lockFile)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// LockedProcess returns the process with the pid in the lock file
// can return os.ErrNotExist or process.ErrorProcessNotRunning
func (l InstanceLock) LockedProcess() (*process.Process, int, error) {
	pid, err := l.LockedPid()
	if err != nil {
		return nil, 0, err
	}
	proc, err := process.NewProcess(int32(pid))
	return proc, pid, err
}

// KillLockedPid kills the process with the pid in the lock file
func (l InstanceLock) KillLockedPid() error {
	proc, pid, err := l.LockedProcess()
	// lockfile or process not found
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, process.ErrorProcessNotRunning) {
		return nil
	}
	if err != nil {
		return err
	}

	err = proc.Terminate()
	if err != nil {
		return err
	}

	proc2, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	_, err = proc2.Wait()
	return err
}
