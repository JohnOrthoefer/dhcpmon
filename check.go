package main

import (
   "log"
)

// check(err, leasesFile)
func check(e error, s string) {
   if e != nil {
      log.Fatalf("%s %v", s, e)
   }
}
