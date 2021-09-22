package main

import (
	"bufio"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"time"
)

const maxEntries = 100

type output struct {
	When time.Time `json:"when"`
	Time int64     `json:"utime"`
	Pipe string    `json:"channel"`
	Text string    `json:"message"`
}

type journalOutput struct {
	Cursor    string `json:"__CURSOR"`
	Timestamp string `json:"__REALTIME_TIMESTAMP"`
	Message   string `json:"MESSAGE"`
	Transport string `json:"_TRANSPORT"`
}

var store *list.List

//
func storeItem(s *list.List, o *output) {
	if s.Len() >= maxEntries {
		s.Remove(s.Front())
	}
	s.PushBack(o)
}

// Dump a list.List to a string.
func getStored(s *list.List) []string {
	var rtn []string

	for e := s.Front(); e != nil; e = e.Next() {
		rtn = append(rtn, (e.Value).(*output).Text)
	}
	return rtn
}

// Just call getStored with right options
func getOutput() []string {
	return getStored(store)
}

func localLogs() ([]byte, error) {
	var rtn []output

	for e := store.Front(); e != nil; e = e.Next() {
		rtn = append(rtn, *(e.Value.(*output)))
	}
	return json.MarshalIndent(rtn, "", "  ")
}

func syslogLogs() ([]byte, error) {
	var rtn []output

	logs := runCmd([]string{
		"/bin/journalctl",
		"--unit=dnsmasq.service",
		"--output=json",
	})

	log.Printf("Log Length %d", logs.Len())

	for e := logs.Front(); e != nil; e = e.Next() {
		var ent journalOutput

		checkWarn(json.Unmarshal([]byte(e.Value.(*output).Text), &ent), "Unmarshal failed")
		t, _ := strconv.ParseInt(ent.Timestamp, 10, 64)
		rtn = append(rtn, output{
			When: time.Unix(t/1000000, t%1000000),
			Time: t / 1000,
			Pipe: ent.Transport,
			Text: ent.Message,
		})
	}

	return json.MarshalIndent(rtn, "", "  ")
}

func GetLogsJson() ([]byte, error) {
	if lookupBool("systemd") {
		return syslogLogs()
	}
	return localLogs()
}

func startScanner(s *list.List, r io.Reader, c string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		o := new(output)
		o.Text = scanner.Text()
		o.When = time.Now()
		o.Time = o.When.Unix()
		o.Pipe = c
		storeItem(s, o)
	}
}

// Start a long running command
func startCmd(cmdln []string) {
	cmd := exec.Command(cmdln[0], cmdln[1:]...)
	//cmdName := path.Base(cmdln[0])

	stdout, err := cmd.StdoutPipe()
	checkFatal(err, "Stdout Pipe")
	stderr, err := cmd.StderrPipe()
	checkFatal(err, "Stderr Pipe")

	log.Printf("Starting '%v'\n", cmdln)
	err = cmd.Start()
	checkFatal(err, fmt.Sprintf("Exec Cmd: %v", cmdln))

	go startScanner(store, stdout, "err")
	go startScanner(store, stderr, "out")
	cmd.Wait()

	log.Fatal("Command Exited")
}

// run a short command
func runCmd(cmdln []string) *list.List {
	cmd := exec.Command(cmdln[0], cmdln[1:]...)

	stdout, err := cmd.StdoutPipe()
	checkFatal(err, "Stdout Pipe")
	stderr, err := cmd.StderrPipe()
	checkFatal(err, "Stderr Pipe")

	log.Printf("Starting '%v'\n", cmdln)
	this := list.New()
	err = cmd.Start()
	checkFatal(err, fmt.Sprintf("Exec Cmd: %v", cmdln))

	go startScanner(this, stdout, "err")
	go startScanner(this, stderr, "out")
	cmd.Wait()

	log.Printf("Done...")
	return this
}

func init() {
	store = list.New()
}

// vim: noai:ts=4:sw=4:set expandtab:
