package runner

import (
	"os"
	"syscall"
)

type Device struct {
	Path  string
	Type  uint32
	Perm  uint32
	Major uint32
	Minor uint32
}

func (d *Device) Mknod() error {
	if _, err := os.Stat(d.Path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	mode := d.Perm | d.Type
	dev := d.Minor | (d.Major << 8)
	return syscall.Mknod(d.Path, mode, int(dev))
}

// allowed devices inside root fs
var defaultDevices = [...]*Device{
	{
		Path:  "/dev/null",
		Type:  syscall.S_IFCHR,
		Perm:  0666,
		Major: 1,
		Minor: 3,
	},
	{
		Path:  "/dev/zero",
		Type:  syscall.S_IFCHR,
		Perm:  0666,
		Major: 1,
		Minor: 5,
	},
	{
		Path:  "/dev/full",
		Type:  syscall.S_IFCHR,
		Perm:  0666,
		Major: 1,
		Minor: 7,
	},
	{
		Path:  "/dev/random",
		Type:  syscall.S_IFCHR,
		Perm:  0666,
		Major: 1,
		Minor: 8,
	},
	{
		Path:  "/dev/urandom",
		Type:  syscall.S_IFCHR,
		Perm:  0666,
		Major: 1,
		Minor: 9,
	},
	{
		Path:  "/dev/tty",
		Type:  syscall.S_IFCHR,
		Perm:  0666,
		Major: 5,
		Minor: 0,
	},
	{
		Path:  "/dev/console",
		Type:  syscall.S_IFCHR,
		Perm:  0666,
		Major: 5,
		Minor: 1,
	},
	{
		Path:  "/dev/ptmx",
		Type:  syscall.S_IFCHR,
		Perm:  0666,
		Major: 5,
		Minor: 2,
	},
}
