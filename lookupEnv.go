package main

import (
	"log"
	"os"
)

func lookupEnv(e string, d string) string {
	val := os.Getenv(e)

	if val == "" {
		val = d
	}

	log.Printf("%s returns %s\n", e, val)

	return val
}

// vim: noai:ts=4:sw=4:set expandtab:
