package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const VERSION = "0.0.1"

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

func dump(host *string, port *int, index *string) {
	uri := &url.URL{
		Scheme: "http",
		Host:   *host + ":" + strconv.Itoa(*port),
		Path:   *index,
	}

	_, err := netClient.Get(uri.String())
	if err != nil {
		log.Fatal("Error when connecting to Elasticsearch:", err)
	}

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
