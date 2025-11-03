package spawn

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/navetacandra/gowle/internal/config"
)

type ChildProcess struct {
	Cmd *exec.Cmd
	mu  sync.Mutex
}

func (p *ChildProcess) Start(cfg *config.GowleConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(cfg.Command) == 0 {
		return nil
	}

	if p.Cmd != nil { // skip when already started
		return errors.New("process already running")
	}

	if p.Cmd != nil && p.Cmd.Process != nil { // skip when already started
		return errors.New("process already running")
	}

	shell := "sh"
	flag := "-c"
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	}

	cmd := exec.Command(shell, flag, cfg.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Dir = cfg.Cwd

	if err := cmd.Start(); err != nil {
		return err
	}

	p.Cmd = cmd
	return nil
}

func (p *ChildProcess) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Cmd == nil || p.Cmd.Process == nil {
		return errors.New("no process running")
	}

	var err error
	if runtime.GOOS == "windows" {
		err = winKill(p.Cmd)
	} else { // assume unix-like
		err = p.Cmd.Process.Signal(os.Interrupt)
	}

	if err != nil {
		fmt.Println(p.Cmd.Process)
		fmt.Println(p.Cmd.Process.Pid)
		return err
	}

	if _, err := p.Cmd.Process.Wait(); err != nil {
		return err
	}

	p.Cmd = nil
	return nil
}

func winHasProcess(pid int) bool {
	cmd := exec.Command("tasklist", "/FI", "PID eq "+strconv.Itoa(pid))
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return false
	}
	return strings.Contains(out.String(), strconv.Itoa(pid))
}

func winKill(proc *exec.Cmd) error {
	if !winHasProcess(proc.Process.Pid) {
		return nil
	}
	cmd := exec.Command("taskkill", "/PID", strconv.Itoa(proc.Process.Pid), "/F", "/T")
	cmd.Dir = proc.Dir
	cmd.Env = proc.Env
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
