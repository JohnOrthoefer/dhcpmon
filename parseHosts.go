package main

import (
	"encoding/json"
	//	"log"
	"os"
	"strings"
    "unicode"
)

const commentChars = "#;"

func stripComment(source string) string {
	if cut := strings.IndexAny(source, commentChars); cut >= 0 {
		return strings.TrimRightFunc(source[:cut], unicode.IsSpace)
	}
	return source
}

func GetHostsJson() ([]byte, error) {
	type HostsJSON struct {
		Ip    string   `json:"ip"`
		Name  string   `json:"name"`
		Alias []string `json:"alias"`
	}

	var hj []HostsJSON

	for _, line := range strings.Split(readHosts(lookupStr("hostsfile")), "\n") {
        tline := stripComment(strings.TrimSpace(line))
		if tline == "" {
			continue
		}
		ent := strings.Fields(tline)
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
