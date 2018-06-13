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

func pickVictim(Sid int, Min int) (Process, error) {

	output := runcmd("ps xhao pid,ppid,pgid,sid,stat")
	scanner := bufio.NewScanner(strings.NewReader(output))
	Plist := make([]Process, 0)
	foundAny := false
	var p Process

	for scanner.Scan() {
		foundAny = true
		i, err := fmt.Sscanf(scanner.Text(), "%d %d %d %d %s", &p.pid, &p.ppid, &p.pgid, &p.sid, &p.stat)
		if err != nil {
			return p, err
		}
		if i == 5 && p.sid == Sid && strings.HasPrefix(p.stat, "R") {
			Plist = append(Plist, p)
		}
	}
	if !foundAny {
		return p, errors.New(fmt.Sprintf("No process with Session ID: %d", Sid))
	}
	if len(Plist) <= Min {
		return p, errors.New(fmt.Sprintf("Minimum threshold %d reached", Min))
	}
	return Plist[rand.Intn(len(Plist))], nil

}
