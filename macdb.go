package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

type OuiEntry struct {
	OUI     string `json:"oui"`
	Private bool   `json:"isPrivate"`
	Company string `json:"companyName"`
	Address string `json:"companyAddress"`
	Country string `json:"contryCode"`
	BlockSz string `json:"assignmentBlockSize"`
	Created string `json:"dateCreated"`
	Updated string `json:"dateUpdated"`
}

type OuiTable map[string]*OuiEntry

var cache OuiTable
var macDB *os.File

func getCompany(i *OuiEntry) string {
    if i == nil {
        return "Unknown"
    }
    return i.Company
}

func lookupMac(mac string) *OuiEntry {
	mac = strings.ToUpper(mac)
    for i := len(mac); i >= 0; i-- {
        val, ok := cache[mac[0:i]]
        if ok {
	        return val
	    }
    }

    if mac[1] == '2' ||
       mac[1] == '6' ||
       mac[1] == 'A' ||
       mac[1] == 'E' {
        return cache["PRIVATE"]
    }

   _, err := macDB.Seek(0,0)
   checkFatal(err, "macDB can not seek")

   scanner := bufio.NewScanner(macDB)
   v:=new(OuiEntry)
   for scanner.Scan() {
      ln := scanner.Bytes()
      json.Unmarshal(ln, &v)
      prefix := strings.ToUpper(v.OUI)
      if prefix == mac[0:len(prefix)] {
         cache[prefix] = v
         return v
      }
    }

	return cache["UNKNOWN"]
}

func initMacs(db string) error {
    var err error
	macDB, err = os.Open(db)
	checkFatal(err, "MAC DB json read")

	return nil
}

func init() {
	cache = make(OuiTable)
    cache["UNKNOWN"] = &OuiEntry{
        OUI:        "00:00:00:00:00:00",
        Private:    false,
        Company:    "UNKNOWN",
        Address:    "UNKNOWN",
        Country:    "",
        BlockSz:    "",
        Created:    "",
        Updated:    "",
    }
    cache["PRIVATE"] = &OuiEntry{
        OUI:        "00:00:00:00:00:00",
        Private:    true,
        Company:    "Privacy Mac",
        Address:    "UNKNOWN",
        Country:    "",
        BlockSz:    "",
        Created:    "",
        Updated:    "",
    }
}

// vim: noai:ts=4:sw=4:set expandtab:
