package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func run() {
	if len(os.Args) < 3 {
		fmt.Println("usage: pocky run <command>")
		os.Exit(0)
	}

	cmd := exec.Command("/proc/self/exe", append([]string{"container"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// hardcoded for now:
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	if cmd.Run() != nil {
		status := cmd.ProcessState.ExitCode()
		os.Exit(status)
	}
}

func container() {
	container := NewContainer()
	container.Run()
}

func printHelp() {
	fmt.Println("Usage: pocker [options] command [arg...]\n")
	fmt.Println("Available commands:")
	fmt.Println("  help - print this help message")
	fmt.Println("  run - executes program inside isolated container")
	fmt.Print("\n")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("pocky: Simple, free and open-source container runner")
		fmt.Println("GitHub: https://github.com/michaldziuba03/pocky\n")
		fmt.Println("usage: pocky run <command>\n")
		os.Exit(0)
	}

	command := os.Args[1]
	switch command {
	case "run":
		run()
		break
	case "help":
		printHelp()
		break
	case "container":
		container()
		break
	default:
		log.Fatal("error: unknown command")
	}
}
