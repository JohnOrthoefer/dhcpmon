package main

import (
   "fmt"
)

const leasesFile="dnsmasq.leases"

func main() {

   readLeases(leasesFile)
   go watchLeases(leasesFile)

   fmt.Printf("Leases %d\n", len(GetLeases()))

   fmt.Printf("Starting HTTP\n")
   startHTTP()
}
