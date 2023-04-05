package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// The CommandExecutor interface defines an
// entity that is able to execute memcached
// commands against a memcached server.
type CommandExecutor interface {
	execute(command string, delimiters []string) []string
	Close()
}

type KeyValue struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Flag   int    `json:"flag,omitempty"`
	Expiry int    `json:"expiry,omitempty"`
	Length int    `json:"length,omitempty"`
}

type Key struct {
	Original string `json:"original,omitempty"`
	Name     string `json:"name,omitempty"`
	Expiry   int    `json:"expiry,omitempty"`
}

type MemcachedCommandExecutor struct {
	connection net.Conn
}

func MemClient(server string) (*memClient, error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}

	return &memClient{
		server: server,
		executor: &MemcachedCommandExecutor{
			connection: conn,
		},
	}, nil
}

type memClient struct {
	server   string
	executor CommandExecutor
}

type Stat struct {
	name  string
	value string
}

func (executor *MemcachedCommandExecutor) execute(command string, responseDelimiters []string) []string {
	fmt.Fprint(executor.connection, command)
	scanner := bufio.NewScanner(executor.connection)
	var result []string

OUTER:
	for scanner.Scan() {
		line := scanner.Text()
		for _, delimiter := range responseDelimiters {
			if line == delimiter {
				break OUTER
			}
		}
		result = append(result, line)
		// if there is no delimiter specified, then the response is just a single line and we should return after
		// reading that first line (e.g. version command)
		if len(responseDelimiters) == 0 {
			break OUTER
		}
		// if there is an error storing the value, then we should break out of the loop
		if line == "ERROR" {
			break OUTER
		}
	}
	return result
}

// Closes the memcached connection
func (executor *MemcachedCommandExecutor) Close() {
	executor.connection.Close()
}

// List all cache keys on the memcached server.
func (client *memClient) ListKeys() []Key {
	var keys []Key
	result := client.executor.execute("stats items\r\n", []string{"END"})

	// identify all slabs and their number of items by parsing the 'stats items' command
	r, _ := regexp.Compile("STAT items:([0-9]*):number ([0-9]*)")
	slabCounts := map[int]int{}
	for _, stat := range result {
		matches := r.FindStringSubmatch(stat)
		if len(matches) == 3 {
			slabId, _ := strconv.Atoi(matches[1])
			slabItemCount, _ := strconv.Atoi(matches[2])
			slabCounts[slabId] = slabItemCount
		}
	}

	// For each slab, dump all items and add each key to the `keys` slice
	r, _ = regexp.Compile(`ITEM\s*?(.*?)\s*?\[\d+\s*?b;\s*?(\d+)\s*?s\]`)
	for slabId, slabCount := range slabCounts {
		command := fmt.Sprintf("stats cachedump %v %v\n", slabId, slabCount)
		commandResult := client.executor.execute(command, []string{"END"})
		for _, item := range commandResult {
			matches := r.FindStringSubmatch(item)
			expiry, _ := strconv.Atoi(matches[2])

			keys = append(keys, Key{
				matches[0],
				strings.TrimSpace(matches[1]),
				expiry,
			})
		}
	}

	return keys
}

// Sets a given cache key on the memcached server to a given value.
func (client *memClient) Set(key string, value string, expiration int) {
	flags := "0"
	command := fmt.Sprintf("set %s %s %d %d\r\n%s\r\n", key, flags, expiration, len(value), value)
	client.executor.execute(command, []string{"STORED"})
}

// Retrieves a given cache key from the memcached server.
// Returns a string array with the value and a boolean indicating
// whether a value was found or not.
func (client *memClient) Get(keyName string) (*KeyValue, bool) {
	command := fmt.Sprintf("get %s\r\n", keyName)
	result := client.executor.execute(command, []string{"END"})

	if len(result) < 2 {
		return &KeyValue{}, false
	}

	keyValue := result[1]

	parts := strings.Split(result[0], " ")
	// ditch the first "VALUE <key> <expiration> <length>" line
	flag, _ := strconv.Atoi(parts[2])
	length, _ := strconv.Atoi(parts[3])

	return &KeyValue{
		keyName,
		keyValue,
		flag,
		0, // Expiry is not returned in response, we will populate it in a different manner
		length,
	}, true
}

// Get the server version.
func (client *memClient) Version() string {
	result := client.executor.execute("version \r\n", []string{})
	if len(result) == 1 {
		return result[0]
	}

	return "UNKNOWN"
}

// Retrieves all server statistics.
func (client *memClient) Stats() []Stat {
	result := client.executor.execute("stats\r\n", []string{"END"})

	var stats []Stat
	for _, stat := range result {
		parts := strings.SplitN(stat, " ", 3)
		stats = append(stats, Stat{parts[1], parts[2]})
	}

	return stats
}

// Retrieves a specific server statistic.
func (client *memClient) Stat(statName string) (Stat, bool) {
	stats := client.Stats()
	for _, stat := range stats {
		if stat.name == statName {
			return stat, true
		}
	}

	return Stat{}, false
}

// Creates a memClient and deals with any errors
// that might occur (e.g. unable to connect to server).
func createClient(addr *string) (*memClient) {
	server := *addr
	client, err := MemClient(server)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to", server)
		os.Exit(1)
	}

	return client
}
