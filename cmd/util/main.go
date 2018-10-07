package main

import (
	"flag"
	"fmt"
	"github.com/op/go-logging"
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

	client.Set("username", "john doe", 900)
	client.Set("age", "3438", 900)
	client.Set("profession", "debugging", 900)
	client.Set("location", "neverland", 900)

	time.Sleep(1000 * time.Millisecond)

	keys := client.ListKeys()
	for _, key := range keys {
		keyValue, _ := client.Get(key)
		fmt.Println(keyValue)
	}
}
