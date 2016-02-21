package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var windows map[string]int64

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func getActiveWindowName() ([]byte, error) {
	c := exec.Command("/usr/bin/osascript", "-e",
		`tell application "System Events"
      set frontmostProcess to name of first process where it is frontmost
    end tell
    return frontmostProcess`)
	return c.CombinedOutput()
}

func printActiveWindow(t time.Time) {
	out, err := getActiveWindowName()
	if err != nil {
		fmt.Println(err)
	}

	windowTitle := strings.Replace(string(out), "\n", "", -1)

	_, ok := windows[windowTitle]
	if !ok {
		windows[windowTitle] = 0
	} else {
		windows[windowTitle]++
	}

	fmt.Printf("%v: %v, total %v seconds\n", t.Format("2006-01-02 15:04:05 MST"), windowTitle, windows[windowTitle])
}

func main() {
	windows = make(map[string]int64)
	doEvery(time.Second, printActiveWindow)
}
