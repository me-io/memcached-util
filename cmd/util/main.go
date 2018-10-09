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
	client    *memClient
	addr      *string
	path      *string
	operation *string

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

	addr = flag.String("addr", `localhost:11211`, "Address to memcached server")
	path = flag.String("filename", "mem_backup.json", "Path to store the output file at")
	operation = flag.String("op", "", "Whether to backup the cache")

	// If the given filename does not have the suffix, add to it
	*path = strings.Trim(*path, "/")
	if !strings.HasSuffix(*path, ".json") {
		*path = *path + ".json"
	}

	flag.Parse()

	client = createClient(addr)
}

// backupCache: Exports the cache into file at the given path
func backupCache() {
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
	Logger.Infof("address %s", *addr)

	if *operation == "backup" {
		backupCache()
	} else if *operation == "restore" {
		restoreCache()
	} else {
		Logger.Error("--op is required with either 'backup' or 'restore'")
		os.Exit(1)
	}
}
