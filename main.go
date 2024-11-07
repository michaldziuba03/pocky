package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

import "pocker/runner"

func run() {
	if len(os.Args) < 3 {
		fmt.Println("usage: pocky run <command>")
		os.Exit(0)
	}

	pid := os.Getpid()
	command := os.Args[2:]
	config := runner.NewConfig(pid, command)
	configJSON, err := json.Marshal(&config)
	if err != nil {
		log.Fatal("error:", err)
	}

	var args = [...]string{"container", string(configJSON)}

	cmd := exec.Command("/proc/self/exe", args[:]...)
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

/*
 * "pocky container <json:config>"
 * NOT PUBLIC for the end user
 */
func container() {
	if len(os.Args) != 3 {
		log.Fatal("error: expected 3 arguments")
	}

	var config runner.Config
	err := json.Unmarshal([]byte(os.Args[2]), &config)
	if err != nil {
		log.Fatal("error: ", err)
	}

	c := runner.NewContainer(&config)
	c.Run()
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
