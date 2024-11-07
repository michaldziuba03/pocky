package runner

type Limits struct {
	MaxPids     int `json:"max_pids"`
	MemoryLimit int `json:"memory_limit"`
}

type Config struct {
	PID        int      `json:"pid"`
	Namespaces int      `json:"namespaces"`
	Command    []string `json:"command"`
	Limits     Limits   `json:"limits"`
}

func NewConfig(pid int, command []string) *Config {
	return &Config{
		PID:        pid,
		Namespaces: 0,
		Command:    command,
		Limits: Limits{
			MaxPids:     -1,
			MemoryLimit: -1,
		},
	}
}
