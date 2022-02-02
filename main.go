package main

import (
	"log"
)

const (
	iniFile = "dhcpmon.ini"
)

var (
	sha1ver   string
	buildTime string
	repoName  string
)

func main() {

	log.Printf("%s: Build %s, Time %s", repoName, sha1ver, buildTime)

	readConfig(iniFile)
	readEnv()

	initMacs(lookupStr("macdbfile"))

	loadTmpl(lookupStr("htmldir"))

    log.Printf("Static = %s", lookupStr("staticfile"))

	leasesFile := lookupStr("leasesfile")
	log.Printf("Leases = %s", leasesFile)
	readLeases(leasesFile)
	go watchLeases(leasesFile)

	if !lookupBool("systemd") {
		go startCmd([]string{
			lookupStr("dnsmasq"),
			"--keep-in-foreground",
			"--conf-dir=/etc/dnsmasq.d,*conf",
		})
	} else {
		//		go startCmd([]string{
		//			"/bin/journalctl",
		//			"--follow",
		//			"--unit=dnsmasq.service",
		//			"--output=json",
		//		})
	}

	log.Printf("Leases %d\n", len(GetLeases()))
	log.Printf("Starting HTTP, %s\n", lookupStr("httplisten"))

	startHTTP()
}

// vim: noai:ts=4:sw=4:set expandtab:
