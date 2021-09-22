package main

import (
	"encoding/json"
	//	"log"
	"os"
	"strings"
)

func GetHostsJson() ([]byte, error) {
	type HostsJSON struct {
		Ip    string   `json:"ip"`
		Name  string   `json:"name"`
		Alias []string `json:"alias"`
	}

	var hj []HostsJSON

	for _, line := range strings.Split(readHosts(lookupStr("hostsfile")), "\n") {
		if line == "" {
			continue
		}
		ent := strings.Fields(line)
		e := HostsJSON{
			Ip:    ent[0],
			Name:  ent[1],
			Alias: ent[2:],
		}
		hj = append(hj, e)
	}
	if len(hj) < 1 {
		return []byte("[]"), nil
	}
	return json.MarshalIndent(hj, "", "  ")
}

func readHosts(f string) string {
	hosts, err := os.ReadFile(f)
	if checkWarn(err, f) {
		return ""
	}

	return string(hosts)
}

// vim: noai:ts=4:sw=4:set expandtab:
