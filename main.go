package main

import (
   "log"
   "fmt"
   "io/ioutil"
)

const leasesFile="dnsmasq.leases"

func check(e error, s string) {
   if e != nil {
      log.Fatalf("%s %v", s, e)
   }
}

func main() {
   leases, err := ioutil.ReadFile(leasesFile)
   check(err, leasesFile)

   d := ParseLeases(string(leases))

   fmt.Printf("Leases %d\n", len(d))

   fmt.Printf("Starting HTTP\n")
   startHTTP()
}
