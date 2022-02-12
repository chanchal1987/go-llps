// llps provides an API for finding and listing processes in a platform-agnostic way.
package llps

import (
	"errors"
	"fmt"

	sdk "github.com/mitchellh/go-ps"
)

var ErrNoProcessFound = errors.New("no process found")

// ErrUnableToFindProcess is returned when a process is cannot be retrieved.
func ErrUnableToFindProcess(err error) error {
	return fmt.Errorf("%s: %w", "unable to find process", err)
}

// Process contains information about a running process.
//
// This is generic to all operating systems supported by github.com/mitchellh/go-ps.
type Process struct {
	// Parent is the information about the parent process.
	Parent *Process

	// PID is the process ID for this process.
	PID int

	// Executable name running this process. This is not a path to the executable.
	Executable string
}

func processMap() (map[int]*Process, error) {
	processes, err := sdk.Processes()
	if err != nil {
		return nil, ErrUnableToFindProcess(err)
	}

	processMap := make(map[int]*Process, len(processes))

	for _, process := range processes {
		ppid := process.PPid()
		if _, ok := processMap[ppid]; !ok && ppid != 0 {
			processMap[ppid] = &Process{Parent: nil, PID: ppid, Executable: ""}
		}

		pid := process.Pid()
		if _, ok := processMap[pid]; !ok {
			processMap[pid] = &Process{Parent: nil, PID: pid, Executable: ""}
		}

		if ppid != 0 {
			processMap[pid].Parent = processMap[ppid]
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
// Thid func will send ErrNoProcessFound if a matching process is not found.
func FindProcess(pid int) (*Process, error) {
	processMap, err := processMap()
	if err != nil {
		return nil, err
	}

	if process, ok := processMap[pid]; ok {
		return process, nil
	}

	return nil, ErrNoProcessFound
}
