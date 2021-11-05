package llps

import sdk "github.com/mitchellh/go-ps"

type Process struct {
	ParentProcess *Process
	PID           int
	Executable    string
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
