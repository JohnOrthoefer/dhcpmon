package main

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
	"strconv"
	"strings"
)

type configTable map[string]string

var configVars configTable

func readConfig(iniFile string) {
	cfg, err := ini.LoadSources(ini.LoadOptions{Insensitive: true}, iniFile)

	if err != nil {
		log.Printf("Skipping %s: %s", iniFile, err)
		return
	}

	defSec := cfg.Section("")
	for k := range configVars {
		if defSec.HasKey(strings.ToLower(k)) {
			v := defSec.Key(strings.ToLower(k)).String()
			configVars[k] = v
		}
	}

}

func readEnv() {
	for k := range configVars {
		v := os.Getenv(strings.ToUpper(k))
		if v != "" {
			configVars[k] = v
		}
	}
}

func lookupStr(s string) string {
	return configVars[strings.ToLower(s)]
}

func lookupBool(s string) bool {
	rtn, err := strconv.ParseBool(configVars[strings.ToLower(s)])
	if err != nil {
		return false
	}
	return rtn
}

func init() {
	configVars = configTable{
		"leasesfile":   "/var/lib/misc/dnsmasq.leases",
		"htmldir":      "/app/html",
		"httplisten":   "127.0.0.1:8067",
		"dnsmasq":      "/usr/sbin/dnsmasq",
		"systemd":      "False",
		"macdbfile":    "/app/macaddress.io-db.json",
		"macdbpreload": "False",
		"nmap":         "/usr/bin/nmap",
		"nmapOpts":     "-oG - -n -F 192.168.12.0/24",
		"hostsfile":    "/var/lib/misc/hosts",
		"httplinks":    "true",
		"httpslinks":   "true",
		"sshlinks":     "true",
        "staticfile":   "/etc/dnsmasq.d/static.conf",
	}
}

// vim: noai:ts=4:sw=4:set expandtab:
