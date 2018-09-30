package main

import (
	"flag"
	"github.com/op/go-logging"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

var (
	host string
	port int

	format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02T15:04:05.999999} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
	)

	// Logger ... Logger Driver
	Logger = logging.MustGetLogger("memcached-util")

	_, filename, _, _ = runtime.Caller(0)
	defaultStaticPath = filepath.Dir(filename) + `/public`
	staticPath        = defaultStaticPath
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

	// Caching
	host = GetEnv(`H`, `0.0.0.0`)
	port, _ = strconv.Atoi(GetEnv(`P`, `5000`))
	staticPath = GetEnv(`STATIC_PATH`, defaultStaticPath)

	flag.Parse()
}

// main ... main function start the server
func main() {

	Logger.Infof("host %s", host)
	Logger.Infof("port %d", port)
	Logger.Infof("Static dir %s", staticPath)

}
