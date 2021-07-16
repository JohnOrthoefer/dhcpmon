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
	}

	log.Printf("Leases %d\n", len(GetLeases()))
	log.Printf("Starting HTTP\n")

	startHTTP()
}

// vim: noai:ts=4:sw=4
