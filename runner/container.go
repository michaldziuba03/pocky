package runner

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Container struct {
	ID     string
	SID    string
	limits CGroup
}

func NewContainer() *Container {
	id := uuid.New().String()
	sid := id[:8] // for display only

	container := &Container{
		ID:     id,
		SID:    sid,
		limits: CGroup{},
	}

	container.limits.InitCGroup(container)
	return container
}

func (c *Container) Run() {
	pid := os.Getpid()
	if pid != 1 {
		log.Fatal("error: cannot run container outside isolated namespace")
	}

	fmt.Printf("Container ID: %s\n\n", c.SID)
	c.setHostname()
	c.setEnvironmentVars()
	c.setRootFS()
	c.mountProc()

	// TODO: pass config, maybe by some JSON in FS...
	cmd := exec.Command(os.Args[2])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if cmd.Run() != nil {
		status := cmd.ProcessState.ExitCode()
		os.Exit(status)
	}

	// cleanup process
	c.unmountProc()
	c.unsetRootFS()
}

func (c *Container) setRootFS() {
	chrootDest := dest + "/alpine"
	oldRoot := chrootDest + "/old"

	err := os.MkdirAll(oldRoot, 0700)
	if err != nil {
		log.Fatal("error: ", err)
	}

	var flags uintptr = syscall.MS_BIND | syscall.MS_REC
	err = syscall.Mount(chrootDest, chrootDest, "", flags, "")
	if err != nil {
		log.Fatal("error: ", err)
	}

	// alternative to syscall.Chroot(chrooDest)
	err = syscall.PivotRoot(chrootDest, oldRoot)
	if err != nil {
		log.Fatal("error: ", err)
	}

	err = syscall.Chdir("/")
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func (c *Container) unsetRootFS() {
	err := syscall.Unmount("/old", syscall.MNT_DETACH)
	if err != nil {
		log.Println("error: ", err)
	}
	err = os.Remove("/old")
	if err != nil {
		log.Println("error: ", err)
	}
}

func (c *Container) mountProc() {
	err := syscall.Mount("proc", "/proc", "proc", 0, "")
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func (c *Container) unmountProc() {
	err := syscall.Unmount("/proc", 0)
	if err != nil {
		log.Println("error: ", err)
	}
}

func (c *Container) setEnvironmentVars() {
	err := os.Setenv("PS1", "\\u@\\h:\\w$ ")
	// tbh, not critical so let's ignore:
	if err != nil {
		return
	}
}

func (c *Container) setHostname() {
	err := syscall.Sethostname([]byte(c.SID))
	if err != nil {
		log.Fatal("error: ", err)
	}
}