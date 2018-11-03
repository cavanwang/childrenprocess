package main

import (
	"flag"
	"fmt"

	"github.com/shirou/gopsutil/process"
	"os"
)

func main() {
	var pid int
	flag.IntVar(&pid, "pid", -1, "pid whose children processes will be listed")
	flag.Parse()
	if pid == -1 {
		fmt.Printf("miss parameter pid\n")
		os.Exit(1)
	}
	pids, err := getGrandson(int32(pid))
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	if len(pids) > 0 {
		fmt.Printf("%d", pids[0])
	}
	for _, p := range pids[1:] {
		fmt.Printf(" %d", p)
	}
	fmt.Println("")
}

func getGrandson(pid int32) ([]int32, error) {
	ps, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("cannot list processes: %v", err)
	}
	found := false
	for i, _ := range ps {
		if ps[i].Pid == pid {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("process not found")
	}

	ret := []int32{pid}
	toList := []int32{pid}
	for {
		cps, err := listProcesses(toList, ps)
		if err != nil {
			return nil, fmt.Errorf("cannot list children: %v", err)
		}
		if len(cps) == 0 {
			break
		}
		ret = append(ret, cps...)
		toList = cps
	}
	return ret, nil
}

func listProcesses(parents []int32, ps []*process.Process) ([]int32, error) {
	ret := []int32{}
	for _, parent := range parents {
		for i, _ := range ps {
			ppid, err := ps[i].Ppid()
			if err != nil {
				return nil, err
			}
			if ppid == parent {
				ret = append(ret, ps[i].Pid)
			}
		}
	}
	return ret, nil
}
