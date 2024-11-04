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
	suid := uid[:8] // short id
	pid := os.Getpid()
	if pid != 1 {
		fmt.Printf("error: Cannot run container outside isolated namespace\n")
		os.Exit(1)
	}

	fmt.Printf("Running [%s] as PID %d\n", os.Args[2], pid)
	fmt.Printf("Container UID: %s\n\n", suid)

	chrootDest := dest + "/alpine"
	oldRoot := chrootDest + "/old"

	err := os.MkdirAll(oldRoot, 0700)
	if err != nil {
		fmt.Printf("error(mkdir): %s\n", err)
		os.Exit(1)
	}

	err = syscall.Mount(chrootDest, chrootDest, "", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	// alternative to syscall.Chroot(chrooDest)
	err = syscall.PivotRoot(chrootDest, oldRoot)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	err = syscall.Chdir("/")
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	err = syscall.Mount("proc", "/proc", "proc", 0, "")
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	err = syscall.Sethostname([]byte(suid))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	err = os.Setenv("PS1", "\\u@\\h:\\w$ ")
	if err != nil {
		return
	}

	cmd := exec.Command(os.Args[2])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		status := cmd.ProcessState.ExitCode()
		os.Exit(status)
	}

	err = syscall.Unmount("/proc", 0)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	err = syscall.Unmount("/old_root", syscall.MNT_DETACH)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	err = os.Remove("/old_root")
	if err != nil {
		fmt.Printf("error: %s\n", err)
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
			status := cmd.ProcessState.ExitCode()
			os.Exit(status)
		}
	case "download":
		Download()
	case "child":
		container()
	default:
		fmt.Println("error: Unknown command", os.Args[1])
		os.Exit(1)
	}
}
