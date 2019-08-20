package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
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

func dial(url string) *http.Response {
	resp, err := netClient.Get(url)
	if err != nil {
		log.Fatal("Error when connecting to Elasticsearch:", err)
	}

	return resp
}

func size(uri *url.URL) int {
	var data map[string]interface{}
	uri.RawQuery = "size=1"
	uri.ForceQuery = true

	resp := dial(uri.String())
	buffer, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(buffer, &data)

	// TODO: retrieve index size
	log.Println(data)
	return 0
}

func dump(uri *url.URL) {
	resp := dial(uri.String())

	buffer, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(buffer))
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

	uri := &url.URL{
		Scheme: "http",
		Host:   *host + ":" + strconv.Itoa(*port),
		Path:   *index + "/_search",
	}

	size(uri)
	dump(uri)
}
