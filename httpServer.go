package main

import (
   "bytes"
   "log"
   "fmt"
   "io/ioutil"
   "net/http"
   "text/template"
)

const baseTMPLFile = "html/bootstrap.tmpl"

type tmplMap map[string]*template.Template
var TMPLs tmplMap

func dhcpTable() string {
   var rtn bytes.Buffer
   TMPLs["Leases"].Execute(&rtn,GetLeases())
   return rtn.String()
}

func handler(w http.ResponseWriter, r *http.Request) {
   var t struct {
      PageTitle string
      PageBody string
   }

   fmt.Printf("Request: %v\n", r.URL.String())
   if r.URL.String() == "/Leases" {
      t.PageTitle = "DHCPmon"
      t.PageBody = dhcpTable()
   }

   TMPLs["Home"].Execute(w, t)
}

func startHTTP() {
   http.HandleFunc("/", handler)
   log.Fatal(http.ListenAndServe(":8000", nil))
}

func init() {
   TMPLs = make(tmplMap)
   in, _ := ioutil.ReadFile(baseTMPLFile)
   TMPLs["Home"] = template.Must(template.New("Home").Parse(string(in)))

   in, _ = ioutil.ReadFile("html/leases.tmpl")
   TMPLs["Leases"] = template.Must(template.New("Leases").Parse(string(in)))
}
