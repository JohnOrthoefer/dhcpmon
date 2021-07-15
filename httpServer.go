package main

import (
   "bytes"
   "log"
   "fmt"
   "io/ioutil"
   "net/http"
   "text/template"
)

type tmplMap map[string]*template.Template
var TMPLs tmplMap
var tmplFiles = map[string]string {
   "Home": "bootstrap.tmpl",
   "Leases": "leases.tmpl",
   "Logs": "logs.tmpl",
}

func dhcpTable() string {
   var rtn bytes.Buffer
   TMPLs["Leases"].Execute(&rtn, GetLeases())
   return rtn.String()
}

func logDisplay() string {
   var rtn bytes.Buffer
   TMPLs["Logs"].Execute(&rtn, getOutput())
   return rtn.String()
}

func handler(w http.ResponseWriter, r *http.Request) {
   var t struct {
      PageTitle string
      PageBody string
   }

   fmt.Printf("Request: %v\n", r.URL.String())
   url := r.URL.String()
   if url == "/Leases" {
      t.PageTitle = "DHCPmon - Leases"
      t.PageBody = dhcpTable()
   } else if url == "/Logs" {
      t.PageTitle = "DHCPmon - Logs"
      t.PageBody = logDisplay()
   }

   TMPLs["Home"].Execute(w, t)
}

func startHTTP() {
   http.HandleFunc("/", handler)
   log.Fatal(http.ListenAndServe(lookupEnv("HTTPLISTEN", ":8067"), nil))
}

func init() {
   htmldir := lookupEnv("HTMLDIR", "/app/html")
   TMPLs = make(tmplMap)

   for key, file := range tmplFiles {
      f:=fmt.Sprintf("%s/%s", htmldir, file)
      log.Printf("Reading...%s\n", f)
      in, err := ioutil.ReadFile(f)
      if err == nil {
         TMPLs[key] = template.Must(template.New("Key").Parse(string(in)))
      } else {
         log.Fatalf("%s: Error %v\n", file, err)
      }
   }
}
