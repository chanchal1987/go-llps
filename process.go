// llps provides an API for finding and listing processes in a platform-agnostic way.
package llps

import sdk "github.com/mitchellh/go-ps"

// Process contains information about a running process.
//
// This is generic to all operating systems supported by github.com/mitchellh/go-ps.
type Process struct {
	// ParentProcess is the information about the parent process.
	ParentProcess *Process

	// Pid is the process ID for this process.
	PID int

	// Executable name running this process. This is not a path to the executable.
	Executable string
}

// String returns a human-friendly string for the process.
func (p *Process) String() string {
	return p.Executable
}

// GoToParent returns nth the parent process of the given process.
func (ps *Process) GoToParent(n int) *Process {
	if n <= 0 {
		return ps
	}

	if ps.ParentProcess == nil {
		return nil
	}

	return ps.ParentProcess.GoToParent(n - 1)
}

// FindExecutable find all the processes matched by the lookup func.
func (ps *Process) FindExecutable(f func(executable string) bool) (processes []*Process) {
	for psCopy := ps; psCopy != nil; psCopy = psCopy.ParentProcess {
		if f(psCopy.Executable) {
			processes = append(processes, psCopy)
		}
	}

	return
}

func processMap() (map[int]*Process, error) {
	processes, err := sdk.Processes()
	if err != nil {
		return nil, err
	}

	processMap := make(map[int]*Process, len(processes))
	for _, process := range processes {
		ppid := process.PPid()
		if _, ok := processMap[ppid]; !ok && ppid != 0 {
			processMap[ppid] = &Process{PID: ppid}
		}

		pid := process.Pid()
		if _, ok := processMap[pid]; !ok {
			processMap[pid] = &Process{PID: pid}
		}

		if ppid != 0 {
			processMap[pid].ParentProcess = processMap[ppid]
		}

		processMap[pid].Executable = process.Executable()
	}

	return processMap, nil
}

// Processes returns all processes.
//
// This of course will be a point-in-time snapshot of when this method was called. Some operating
// systems don't provide snapshot capability of the process table, in which case the process table
// returned might contain ephemeral entities that happened to be running when this was called.
func Processes() ([]*Process, error) {
	processMap, err := processMap()
	if err != nil {
		return nil, err
	}

	processes := make([]*Process, 0, len(processMap))
	for _, process := range processMap {
		processes = append(processes, process)
	}

	return processes, nil
}

// FindProcess looks up a single process by pid.
//
// Process will be nil and error will be nil if a matching process is not found.
func FindProcess(pid int) (*Process, error) {
	processMap, err := processMap()
	if err != nil {
		return nil, err
	}

	if process, ok := processMap[pid]; ok {
		return process, nil
	}

	return nil, nil
}
