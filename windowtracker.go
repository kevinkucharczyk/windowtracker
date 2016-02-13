package main

import (
	"fmt"
	"os/exec"
	"time"
)

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
	fmt.Printf("%v: %v", t.Format("2006-01-02 15:04:05 MST"), string(out))
}

func main() {
	doEvery(time.Second, printActiveWindow)
}
