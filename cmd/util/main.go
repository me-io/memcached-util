package main

import (
	"encoding/json"
	"flag"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var (
	client  *memClient
	host    *string
	port    *string
	path    *string
	backup  *bool
	restore *bool

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

	host = flag.String("host", `0.0.0.0`, "Memcached hostname")
	port = flag.String("port", "11211", "Memcached port")
	path = flag.String("name", "output.json", "Path to store the output file at")
	backup = flag.Bool("backup", false, "Whether to backup the cache")
	restore = flag.Bool("restore", false, "Whether to restore the cache")

	// If the given filename does not have the suffix, add to it
	*path = strings.Trim(*path, "/")
	if !strings.HasSuffix(*path, ".json") {
		*path = *path + ".json"
	}

	flag.Parse()

	client = createClient(host, port)
}

// backupCache: Exports the cache into file at the given path
func backupCache() {
	client.Set("username", "john doe", 60)
	client.Set("age", "3438", 80)
	client.Set("profession", "debugging", 90)
	client.Set("location", "neverland", 602)

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
		keyValue, _ := client.Get(key.Name)
		keyValue.Expiry = key.Expiry
		cachedData = append(cachedData, *keyValue)
	}

	cachedJson, _ := json.Marshal(cachedData)
	ioutil.WriteFile(*path, cachedJson, 0644)
	Logger.Infof("Output file successfully generated at: %s", *path)
}

// restoreCache: Checks for the existence of the given file and
// restores the data back to memcached
func restoreCache() {
	if _, err := os.Stat(*path); os.IsNotExist(err) {
		Logger.Errorf("File %s does not exist or is not readable", *path)
		os.Exit(1)
	}

	cachedData, _ := ioutil.ReadFile(*path)
	var keyValues []KeyValue
	err := json.Unmarshal(cachedData, &keyValues)
	if err != nil {
		panic(err)
	}

	for _, keyValue := range keyValues {
		expiryTime := time.Unix(int64(keyValue.Expiry), 0)
		currentTime := time.Now()

		duration := expiryTime.Sub(currentTime)
		expirySeconds := int(duration.Seconds())
		if expirySeconds <= 0 {
			Logger.Warningf("Key %s already expired, skipping ..", keyValue.Name)
		} else {
			Logger.Infof("Restoring value for %s. Expires in %d seconds", keyValue.Name, expirySeconds)
		}

		client.Set(keyValue.Name, keyValue.Value, int(expirySeconds))
	}
}

// main: Validates the arguments and processes backup or restore
func main() {
	Logger.Infof("host %s", *host)
	Logger.Infof("port %s", *port)

	// If both the options are given or none of the options are given
	if (*backup && *restore) || (!*backup && !*restore) {
		Logger.Error("Exactly one option --backup or --restore is required")
		os.Exit(1)
	}

	if *backup {
		backupCache()
	} else if *restore {
		restoreCache()
	}
}
