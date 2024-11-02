package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func container() {
	fmt.Printf("Running [%s] as PID %d\n", os.Args[2], os.Getpid())

	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		fmt.Println("Error mounting /proc:", err)
		os.Exit(1)
	}
	defer syscall.Unmount("/proc", 0)

	cmd := exec.Command(os.Args[2])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run container.go run <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC, // Nowa przestrzeń mount
		}

		if err := cmd.Run(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "child":
		container()
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}
