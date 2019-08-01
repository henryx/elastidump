package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/url"
	"strconv"
)

const VERSION = "0.0.1"

func dump(host *string, port *int, index *string) {
	uri := &url.URL{
		Scheme: "http",
		Host:   *host + ":" + strconv.Itoa(*port),
		Path:   *index,
	}

	fmt.Println("GET", uri)
}

func main() {
	host := kingpin.Flag("host", "Set the Elasticsearch host (default localhost)").
		Short('H').
		Default("localhost").
		String()
	port := kingpin.Flag("port", "Set the Elasticsearch port (default 9200)").
		Short('P').
		Default("9200").
		Int()
	index := kingpin.Arg("index", "Index to dump").
		Required().
		String()

	kingpin.HelpFlag.Short('h')
	kingpin.Version(VERSION)
	kingpin.Parse()

	dump(host, port, index)
}
