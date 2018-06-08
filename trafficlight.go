package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"
)

type SwapState struct {
	isSwapping bool
	count      int
}

var stoplist []Process

func readSwapCount() int {
	file, err := os.Open("/proc/vmstat")
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer file.Close()
	for {
		var swapcount int
		var n int
		n, err = fmt.Fscanf(file, "pswpout %d", &swapcount)
		if n == 1 {
			return swapcount
		}
	}

	return -1
}

func main() {
	var swap SwapState
	x := readSwapCount()
	swap.isSwapping = false
	swap.count = x

	Sid := flag.Int("sid", -1, "Session ID of processess to manage")
	Min := flag.Int("min", 1, "Minimum amount of compiles to run at the same time")
	flag.Parse()

	for {
		x := readSwapCount()
		if x > swap.count {
			swap.isSwapping = true
			swap.count = x
			process, err := pickVictim(*Sid, *Min)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Stopping:", process.pid)
				syscall.Kill(process.pid, syscall.SIGSTOP)
				stoplist = append(stoplist, process)
			}
		} else {
			swap.isSwapping = false
			if len(stoplist) > 0 {
				reanimate := stoplist[len(stoplist)-1]
				fmt.Println("Re-animate:", reanimate.pid)
				syscall.Kill(reanimate.pid, syscall.SIGCONT)
				stoplist = stoplist[:len(stoplist)-1]
			}
		}
		time.Sleep(5 * time.Second)

	}
}
