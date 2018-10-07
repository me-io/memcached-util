package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"time"
)

var (
	host *string
	port *string

	format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02T15:04:05.999999} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
	)

	// Logger ... Logger Driver
	Logger = logging.MustGetLogger("memcached-util")
)

// init ... init function of the server
func init() {
	// Logging
	backendStderr := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatted := logging.NewBackendFormatter(backendStderr, format)
	// Only DEBUG and more severe messages should be sent to backend1
	backendLevelFormatted := logging.AddModuleLevel(backendFormatted)
	backendLevelFormatted.SetLevel(logging.DEBUG, "")
	// Set the backend to be used.
	logging.SetBackend(backendLevelFormatted)

	host = flag.String("H", `0.0.0.0`, "Memcached hostname")
	port = flag.String("P", "11211", "Memcached port")

	flag.Parse()
}

// main ... main function start the server
func main() {
	Logger.Infof("host %s", *host)
	Logger.Infof("port %d", *port)

	// connect to memcached server
	client := createClient(host, port)

	client.Set("username", "john doe", 20)
	client.Set("age", "3438", 20)
	client.Set("profession", "debugging", 20)
	client.Set("location", "neverland", 20)

	time.Sleep(1000 * time.Millisecond)

	var cachedData []KeyValue
	keys := client.ListKeys()
	foundCount := len(keys)

	Logger.Infof("%d values found in the storage", foundCount)
	if foundCount == 0 {
		Logger.Infof("No records to publish")
		os.Exit(0)
	}

	for _, key := range keys {
		keyValue, _ := client.Get(key)
		cachedData = append(cachedData, *keyValue)
	}

	cachedJson, _ := json.Marshal(cachedData)
	ioutil.WriteFile("output.json", cachedJson, 0644)
	fmt.Printf("%+v", cachedData)
}
