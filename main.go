package main

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"syscall"
)

func container() {
	uid := uuid.New().String()
	pid := os.Getpid()
	if pid != 1 {
		panic("Cannot run container outside isolated namespace")
	}

	fmt.Printf("Running [%s] as PID %d\n", os.Args[2], pid)
	fmt.Printf("Container id UID: %s\n\n", uid[:8])

	err := syscall.Sethostname([]byte(uid[:8]))
	if err != nil {
		fmt.Printf("Error setting hostname: %s\n", err)
		os.Exit(1)
	}

	err = syscall.Mount("proc", "/proc", "proc", 0, "")
	if err != nil {
		fmt.Printf("Error mounting proc: %s\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[2])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		status := cmd.ProcessState.ExitCode()
		os.Exit(status)
	}

	err = syscall.Unmount("/proc", 0)
	if err != nil {
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: pocky run <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
			Unshareflags: syscall.CLONE_NEWNS,
		}

		if err := cmd.Run(); err != nil {
			fmt.Println("Error:", err)
			status := cmd.ProcessState.ExitCode()
			os.Exit(status)
		}
	case "child":
		container()
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}
