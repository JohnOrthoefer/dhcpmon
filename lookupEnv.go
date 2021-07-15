package main

import (
   "os"
   "log"
)

func lookupEnv(e string, d string) string {
   val := os.Getenv(e)

   if val == "" {
      val = d
   }

   log.Printf("%s returns %s\n", e, val)

   return val
}
