# Trafficlight

A quick hack to stop/continue processes with specific SID when system starts paging

## Problem statement

You have a machine with lots of cores but limited memory. You want to compile with all cores, but occasionally you'll run our of memory and swapping grinds the build to halt. You could:

1. Buy more memory.

   Too obvious solution.

2. Reduce parallellism to limit concurrent memory consumption

   For example if at peak compilers take 1GB, you have 4GB, build with -j5. Not a good solution if you
   have 24 cores and most of the time builds would fit in RAM.

3. Use this tool to stop too many processess trying to page at the same time

## Usage

Use ps to find the SID, the session ID covering your compilers, and run trafficlight on it
Set --min to the minumum amount of compiles to run at the same time (default 1)

```
ps xhao pid,ppid,pgid,sid,stat,comm|grep ninja
26484 26344 26104 25930 S+   ninja
sudo ./trafficlight --min 4 --sid 25930
```

## Authors

(C) Riku Voipio 2018

