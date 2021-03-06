/*
   Copyright (C) 2019 Enrico Bianchi (enrico.bianchi@gmail.com)
   Project       elastidump
   Description   An Elasticsearch dump index
   License       GPL version 2 (see LICENSE for details)
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const VERSION = "1.0.0"

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

type hitData struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type"`
	Id     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source,string,omitempty"`
}

type shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type hits struct {
	Total    int       `json:"total"`
	MaxScore float64   `json:"max_score"`
	Hits     []hitData `json:"hits,omitempty"`
}

type response struct {
	ScrollId string `json:"_scroll_id,omitempty"`
	Took     int    `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   shards `json:"_shards"`
	Hits     hits   `json:"hits"`
}

func dialGet(url string) *http.Response {
	resp, err := netClient.Get(url)
	if err != nil {
		log.Fatal("Error when connecting to Elasticsearch:", err)
	}

	return resp
}

func dialPost(url string, r io.Reader) *http.Response {
	resp, err := netClient.Post(url, "application/json; charset=utf-8", r)
	if err != nil {
		log.Fatal("Error when connecting to Elasticsearch:", err)
	}

	return resp
}

func size(uri *url.URL) int {
	var data response

	resp := dialGet(uri.String())
	buffer, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(buffer, &data)

	return data.Hits.Total
}

func extractData(uri *url.URL) []hitData {
	var data response

	resp := dialGet(uri.String())

	buffer, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Fatal("Query error: ", string(buffer))
	}

	_ = json.Unmarshal(buffer, &data)

	return data.Hits.Hits
}

func extractScroll(uri *url.URL, query string) (string, []hitData) {
	var data response

	resp := dialPost(uri.String(), bytes.NewBuffer([]byte(query)))

	buffer, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Fatal("Query error: ", string(buffer))
	}

	_ = json.Unmarshal(buffer, &data)

	return data.ScrollId, data.Hits.Hits
}

func write(data []hitData) {
	for _, hit := range data {
		document, _ := json.Marshal(hit.Source)
		fmt.Printf("%s\n", string(document))
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

	uri := &url.URL{
		Scheme:     "http",
		Host:       *host + ":" + strconv.Itoa(*port),
		Path:       *index + "/_search",
		ForceQuery: false,
	}

	uri.RawQuery = "size=1"
	total := size(uri)

	if total <= 10000 {
		uri.RawQuery = "size=" + strconv.Itoa(total)
		data := extractData(uri)

		write(data)
	} else {
		uri.RawQuery = "scroll=10m"

		query := `{"size": 100, "query": {"match_all": {}}}`
		for {
			id, data := extractScroll(uri, query)
			if len(data) == 0 {
				break
			} else {
				uri.RawQuery = ""
				uri.Path = "_search/scroll"

				query = fmt.Sprintf(`{"scroll": "10m", "scroll_id": "%s"}`, id)
			}

			write(data)
		}
	}
}
