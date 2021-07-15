package main

import (
   "fmt"
)

const leasesFile="dnsmasq.leases"

func main() {

   leasesFile:=lookupEnv("LEASESFILE", "/var/lib/misc/dnsmasq.leases")
   dnsmasqExec:=lookupEnv("DNSMASQ", "/usr/sbin/dnsmasq")

   readLeases(leasesFile)
   go watchLeases(leasesFile)

   go startCmd([]string{
      dnsmasqExec,
      "--keep-in-foreground",
      "--conf-dir=/etc/dnsmasq.d,*conf",
   })

   fmt.Printf("Leases %d\n", len(GetLeases()))

   fmt.Printf("Starting HTTP\n")

   startHTTP()
}
