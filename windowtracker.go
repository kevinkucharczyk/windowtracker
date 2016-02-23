package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

var windows map[string]*Entry
var signalChan chan os.Signal

type Entry struct {
	TotalSeconds int64
	LastActiveTime time.Time
	Windows map[string]*Entry
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func getActiveProcessName() ([]byte, error) {
	c := exec.Command("/usr/bin/osascript", "-e",
		`tell application "System Events"
      set frontmostProcess to name of first process where it is frontmost
    end tell
    return frontmostProcess`)
	return c.CombinedOutput()
}

func getFrontWindowName(application string) ([]byte, error) {
	c := exec.Command("/usr/bin/osascript", "-e",
		`tell application "` + application + `"
      set windowName to name of front window
    end tell
    return windowName`)
	return c.CombinedOutput()
}

func printActiveWindow(t time.Time) {
	activeProcess, err := getActiveProcessName()
	if err != nil {
		fmt.Println(err)
	}

	activeProcessTitle := strings.Replace(string(activeProcess), "\n", "", -1)

	entry, ok := windows[activeProcessTitle]
	if !ok {
		windows[activeProcessTitle] = &Entry{TotalSeconds: 0, LastActiveTime: t, Windows: make(map[string]*Entry)}
		entry = windows[activeProcessTitle]
	} else {
		entry.TotalSeconds++
		entry.LastActiveTime = t
	}

	fmt.Printf("%v: %v, total %v seconds\n", t.Format("2006-01-02 15:04:05 MST"), activeProcessTitle, entry.TotalSeconds)

	frontWindow, err := getFrontWindowName(activeProcessTitle)

	if err == nil {
		frontWindowTitle := strings.Replace(string(frontWindow), "\n", "", -1)

		windowEntry, ok := entry.Windows[frontWindowTitle]
		if !ok {
			entry.Windows[frontWindowTitle] = &Entry{TotalSeconds: 0, LastActiveTime: t}
			windowEntry = entry.Windows[frontWindowTitle]
		} else {
			windowEntry.TotalSeconds++
			windowEntry.LastActiveTime = t
		}

		fmt.Printf("%v: %v, total %v seconds\n", t.Format("2006-01-02 15:04:05 MST"), frontWindowTitle, windowEntry.TotalSeconds)
	}
}

func beforeExit() {
	<-signalChan
	fmt.Println("Saving history to file...")
	saveJsonFile(windows, "./history.json")
	os.Exit(1)
}

func saveJsonFile(v interface{}, path string) {
	fo, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	e := json.NewEncoder(fo)
	if err := e.Encode(v); err != nil {
		panic(err)
	}
}

func main() {
	windows = make(map[string]*Entry)
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go beforeExit()
	doEvery(time.Second, printActiveWindow)
}
