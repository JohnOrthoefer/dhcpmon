package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

type OuiJSON struct {
	OUI     string `json:"oui"`
	Private bool   `json:"isPrivate"`
	Company string `json:"companyName"`
	Address string `json:"companyAddress"`
	Country string `json:"contryCode"`
	BlockSz string `json:"assignmentBlockSize"`
	Created string `json:"dateCreated"`
	Updated string `json:"dateUpdated"`
}

type OuiTable map[string]string

var cache OuiTable

func lookupMac(mac string) string {
	mac = strings.ToUpper(mac[0:8])
	val, ok := cache[mac]
	if ok {
		return val
	}

	return "MAC not found"
}

func initMacs(db string) (int, error) {
	macDB, err := os.Open(db)
	checkFatal(err, "MAC DB json read")

	defer macDB.Close()

	scanner := bufio.NewScanner(macDB)
	for scanner.Scan() {
		var v OuiJSON
		ln := scanner.Bytes()
		json.Unmarshal(ln, &v)
		cache[strings.ToUpper(v.OUI)] = v.Company
	}

	return len(cache), nil
}

func init() {
	cache = make(OuiTable)
}
