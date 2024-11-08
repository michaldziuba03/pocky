package runner

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const CGroupUnified = "/sys/fs/cgroup/unified"
const CGroupPath = "/sys/fs/cgroup"
const CGroupMemory = "memory"
const CGroupCpu = "cpu"
const CGroupPids = "pids"

type PidsCGroup struct {
	path   string
	pidSet bool
}

func tryCreateDir(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func (ctrl *PidsCGroup) SetMaxPids(limit int) {
	tryCreateDir(ctrl.path)
	str := strconv.Itoa(limit)
	path := filepath.Join(ctrl.path, "pids.max")

	err := os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		log.Fatal("error: ", err)
	}

	ctrl.TrySetProcs()
}

func (ctrl *PidsCGroup) TrySetProcs() {
	if ctrl.pidSet {
		return
	}

	str := strconv.Itoa(os.Getpid())
	path := filepath.Join(ctrl.path, "cgroup.procs")
	err := os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		log.Fatal("error: ", err)
	}

	ctrl.pidSet = true
}

type MemoryCGroup struct {
	path   string
	pidSet bool
}

func (ctrl *MemoryCGroup) SetMemoryLimit(limit int64) {
	tryCreateDir(ctrl.path)
	str := strconv.FormatInt(limit, 10)
	path := filepath.Join(ctrl.path, "memory.limit_in_bytes")

	err := os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		log.Fatal("error: ", err)
	}

	ctrl.TrySetProcs()
}

func (ctrl *MemoryCGroup) TrySetProcs() {
	if ctrl.pidSet {
		return
	}

	str := strconv.Itoa(os.Getpid())
	path := filepath.Join(ctrl.path, "cgroup.procs")
	err := os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		log.Fatal("error: ", err)
	}

	ctrl.pidSet = true
}

type CGroups struct {
	Pids   PidsCGroup
	Memory MemoryCGroup
}

func NewCGroups(c *Container) *CGroups {
	return &CGroups{
		Pids: PidsCGroup{
			path:   CGroupPathV1(CGroupPids, c),
			pidSet: false,
		},
		Memory: MemoryCGroup{
			path:   CGroupPathV1(CGroupMemory, c),
			pidSet: false,
		},
	}
}

func (cg *CGroups) SetCGroups(c *Container) {
	limits := c.config.Limits

	if limits.MemoryLimit != -1 {
		cg.Memory.SetMemoryLimit(limits.MemoryLimit)
	}

	if limits.MaxPids != -1 {
		cg.Pids.SetMaxPids(limits.MaxPids)
	}
}

func CGroupPathV1(controller string, c *Container) string {
	return filepath.Join(CGroupPath, controller, "pocky", c.ID)
}
