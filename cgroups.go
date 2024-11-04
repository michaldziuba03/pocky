package main

import (
	"os"
	"path/filepath"
	"strconv"
)

func SetCgroups(uid string) {
	const path = "/sys/fs/cgroup"
	memory := filepath.Join(path, "memory")
	cgroupMem := filepath.Join(memory, uid)

	if os.MkdirAll(cgroupMem, 0755) != nil {
		panic("Failed to create cgroup memory")
	}

	bytesLimit := filepath.Join(cgroupMem, "memory.limit_in_bytes")
	limit := 104857600 // set 100 MB limit

	if os.WriteFile(bytesLimit, []byte(strconv.Itoa(limit)), 0644) != nil {
		panic("Failed set memory limit in cgroup")
	}

	oomControl := filepath.Join(cgroupMem, "memory.oom_control")
	if os.WriteFile(oomControl, []byte("1"), 0644) != nil {
		panic("Failed set memory oom control in cgroup")
	}

	onRelease := filepath.Join(cgroupMem, "notify_on_release")
	if os.WriteFile(onRelease, []byte("1"), 0644) != nil {
		panic("Failed set notify_on_release in memory cgroup")
	}

	pid := os.Getpid()
	cgroupProcs := filepath.Join(cgroupMem, "cgroup.procs")
	if os.WriteFile(cgroupProcs, []byte(strconv.Itoa(pid)), 0644) != nil {
		panic("Failed set memory limit in cgroup")
	}
}
