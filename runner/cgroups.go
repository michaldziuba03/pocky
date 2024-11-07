package runner

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const CGroupPath = "/sys/fs/cgroup/pocky"
const CGroupMemory = "memory"
const CGroupCpu = "cpu"
const CGroupPids = "pids"

type CGroup struct {
	MemoryLimit  int64
	CpuLimit     int64
	ProcessLimit int64
	container    *Container
}

func (cg *CGroup) InitCGroup(container *Container) {
	cg.container = container
	cg.CpuLimit = -1
	cg.ProcessLimit = -1
	cg.MemoryLimit = -1

	err := os.MkdirAll(CGroupPath, 0755)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func (cg *CGroup) SetMemoryLimit(limit int64) {
	cg.MemoryLimit = limit
	cg.setCGroup(CGroupMemory, "limit_in_bytes", strconv.FormatInt(limit, 10))
}

func (cg *CGroup) SetProcessLimit(limit int64) {
	cg.ProcessLimit = limit
	cg.setCGroup(CGroupPids, "max", strconv.FormatInt(limit, 10))
}

func (cg *CGroup) pathTo(cgroupName string) string {
	return filepath.Join(CGroupPath, cgroupName, cg.container.ID)
}

func (cg *CGroup) setCGroup(cgroupName string, param string, value string) {
	path := cg.pathTo(cgroupName)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatal("error: ", err)
	}

	filename := cgroupName + "." + param
	file := filepath.Join(path, filename)
	err = os.WriteFile(file, []byte(value), 0644)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func (cg *CGroup) addPID(cgroupName string) {
	path := cg.pathTo(cgroupName)
	pid := os.Getpid()
	procs := filepath.Join(path, "cgroup.procs")
	err := os.WriteFile(procs, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		log.Fatal("error: ", err)
	}
}
