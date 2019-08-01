package main

import "gopkg.in/alecthomas/kingpin.v2"

const VERSION = "0.0.1"

func main() {
	kingpin.Version(VERSION)
	kingpin.Parse()
}
