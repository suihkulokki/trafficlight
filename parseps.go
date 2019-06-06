package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
)

type Process struct {
	pid  int
	ppid int
	pgid int
	sid  int
	stat string
}

func runcmd(command string) string {

	out, err := exec.Command("/bin/sh", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}

// Lower min if low on swap
func lowerMinima(min int) (int, error) {
	var swapFree int
	output := runcmd("sed -n /SwapFree:/p /proc/meminfo")
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		i, err := fmt.Sscanf(scanner.Text(), "SwapFree: %d kB", &swapFree)
		if err != nil {
			return min, err
		}
		if i == 1 && swapFree < 1000000 {
			return 1, nil
		}
	}
	return min, nil
}

func pickVictim(Sid int, Min int) (Process, error) {

	output := runcmd("ps xhao pid,ppid,pgid,sid,stat")
	scanner := bufio.NewScanner(strings.NewReader(output))
	Rlist := make([]Process, 0)
	Dlist := make([]Process, 0)
	var p Process
	Min, err := lowerMinima(Min)

	if err != nil {
		return p, err
	}
	foundAny := false

	for scanner.Scan() {
		foundAny = true
		i, err := fmt.Sscanf(scanner.Text(), "%d %d %d %d %s", &p.pid, &p.ppid, &p.pgid, &p.sid, &p.stat)
		if err != nil {
			return p, err
		}
		if i == 5 && p.sid == Sid && strings.HasPrefix(p.stat, "D") {
			Dlist = append(Dlist, p)
		}
		if i == 5 && p.sid == Sid && strings.HasPrefix(p.stat, "R") {
			Rlist = append(Rlist, p)
		}
	}
	if !foundAny {
		return p, errors.New(fmt.Sprintf("No process with Session ID: %d", Sid))
	}
	if (len(Dlist) + len(Rlist)) <= Min {
		return p, errors.New(fmt.Sprintf("Minimum threshold %d reached", Min))
	}
	if len(Dlist) > 0 {
		return Dlist[rand.Intn(len(Dlist))], nil
	}
	return Rlist[rand.Intn(len(Rlist))], nil

}
