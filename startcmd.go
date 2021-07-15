package main

import (
   "fmt"
   "log"
   "io"
   "bufio"
   "os/exec"
   "time"
   "container/list"
)

const maxEntries=100

type output struct {
   when time.Time
   pipe string
   text string
}

var store *list.List

func storeItem(o *output) {
   if store.Len() >= maxEntries {
      store.Remove(store.Front())
   }
   store.PushBack(o)
}

func getOutput()[]string {
   var rtn []string

   for e:= store.Front(); e != nil; e = e.Next() {
      fmt.Printf("%v\n", (e.Value).(*output).text)
      rtn = append(rtn, (e.Value).(*output).text)
   }
   return rtn
}

func startScanner(r io.Reader, c string) {
   scanner := bufio.NewScanner(r)
   for scanner.Scan() {
      o := new(output)
      o.text = scanner.Text()
      o.when = time.Now()
      o.pipe = c
      storeItem(o)
      fmt.Printf("%s: %v\n", o.pipe, o.text)
   }
}

func startCmd(cmdln []string) {
   cmd := exec.Command(cmdln[0], cmdln[1:]...)
   //cmdName := path.Base(cmdln[0])

   stdout, err := cmd.StdoutPipe()
   check(err, "Stdout Pipe")
   stderr, err := cmd.StderrPipe()
   check(err, "Stderr Pipe")

   log.Printf("Starting '%v'\n", cmdln)
   err = cmd.Start()
   check(err, fmt.Sprintf("Exec Cmd: %v", cmdln))

   go startScanner(stdout, "err")
   go startScanner(stderr, "out")
   cmd.Wait()

   log.Fatal("Command Exited")
}

func init() {
   store = list.New()
}
