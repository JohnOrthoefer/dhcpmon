package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type tmplMap map[string]*template.Template
type HTMLOptions struct {
	EnableHTTPLinks  bool
	EnableHTTPSLinks bool
	EnableSSHLinks   bool
	EnableNetworkTags   bool
	EnableEdit   bool
}

var TMPLs tmplMap
var tmplFiles = map[string]string{
	"Home":   "bootstrap.tmpl",
	"Leases": "leases.tmpl",
	"Logs":   "logs.tmpl",
	"Splash": "splash.tmpl",
	"Hosts":  "hosts.tmpl",
}

func dhcpTable() string {
	var rtn bytes.Buffer
	TMPLs["Leases"].Execute(&rtn, HTMLOptions{
		lookupBool("httplinks"),
		lookupBool("httpslinks"),
		lookupBool("sshlinks"),
		lookupBool("networktags"),
		lookupBool("edit"),
	})
	return rtn.String()
}

func logDisplay() string {
	var rtn bytes.Buffer
	TMPLs["Logs"].Execute(&rtn, getOutput())
	return rtn.String()
}

func hostsDisplay() string {
	var rtn bytes.Buffer
	TMPLs["Hosts"].Execute(&rtn, getOutput())
	return rtn.String()
}

func displaySplash() string {
	var rtn bytes.Buffer
	TMPLs["Splash"].Execute(&rtn, "")
	return rtn.String()
}

func handler(w http.ResponseWriter, r *http.Request) {
	var t struct {
		PageTitle string
		PageBody  string
	}

	urlPath := r.URL.String()
	urlQuery := r.URL.Query()
	log.Printf("Request from %s: url=%v query=%v\n", r.Host, urlPath, urlQuery)
    log.Printf("Body: %v\n", r.FormValue("data"))

	if v, found := urlQuery["api"]; found {
		log.Printf("API = %v\n", v[0])
		switch v[0] {
		case "logs.json":
			logs, _ := GetLogsJson()
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte("{\"data\":" + string(logs) + "}"))
			return
		case "leases.json":
			leases, _ := GetLeasesJson()
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte("{\"data\":" + string(leases) + "}"))
			return
		case "hosts.json":
			hosts, _ := GetHostsJson()
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte("{\"data\":" + string(hosts) + "}"))
			return
		}
	} else if v, found := urlQuery["p"]; found {
		log.Printf("Page = %v\n", v[0])
		switch v[0] {
		case "Leases":
			t.PageTitle = "DHCPmon - Leases"
			t.PageBody = dhcpTable()
		case "Logs":
			t.PageTitle = "DHCPmon - Logs"
			t.PageBody = logDisplay()
		case "Hosts":
			t.PageTitle = "DHCPmon - Hosts"
			t.PageBody = hostsDisplay()
		default:
			t.PageTitle = "DHCPmon - Splash"
			t.PageBody = displaySplash()
		}
	}

	TMPLs["Home"].Execute(w, t)
}

func startHTTP() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(lookupStr("httplisten"), nil))
}

func loadTmpl(htmldir string) {
	TMPLs = make(tmplMap)

	for key, file := range tmplFiles {
		f := fmt.Sprintf("%s/%s", htmldir, file)
		log.Printf("Reading...%s\n", f)
		in, err := ioutil.ReadFile(f)
		if err == nil {
			TMPLs[key] = template.Must(template.New("Key").Parse(string(in)))
		} else {
			log.Fatalf("%s: Error %v\n", file, err)
		}
	}
}

// vim: noai:ts=4:sw=4:set expandtab:
