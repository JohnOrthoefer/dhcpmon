package main

import (
	"log"
)

// check(err, leasesFile)
func checkFatal(e error, s string) {
	if e != nil {
		log.Fatalf("%s %v", s, e)
	}
}

func checkWarn(e error, s string) {
	if e != nil {
		log.Printf("%s %v", s, e)
	}
}
