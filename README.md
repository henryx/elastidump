# Elastidump

An Elasticsearch index dump tool

## Description

Elastidump extract all documents saved in an Elasticsearch index and return them to the output. Is useful for a textual
backup of the index data

## Install

`go install github.com/henryx/elastidump`

## Usage

Simple usage is `elastidump twitter` to dump content of index `twitter` to the standard output.

Here is a list of the command line arguments:

```
usage: elastidump [<flags>] <index>

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and
                          --help-man).
  -H, --host="localhost"  Set the Elasticsearch host (default localhost)
  -P, --port=9200         Set the Elasticsearch port (default 9200)
      --version           Show application version.

Args:
  <index>  Index to dump
```