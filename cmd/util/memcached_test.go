// +build !integration

package main

import (
	"fmt"
	"reflect"
	"testing"
)

// Test Helpers
type MockedCommandExecutor struct {
	t                *testing.T
	executedCommands []string
	returnValues     map[string][]string
	closed           bool
}

func (executor *MockedCommandExecutor) execute(command string, responseDelimiters []string) []string {
	executor.executedCommands = append(executor.executedCommands, command)
	returnVal, ok := executor.returnValues[command]
	if ok {
		return returnVal
	}
	return []string{}
}

func (executor *MockedCommandExecutor) Close() {
}

func (executor *MockedCommandExecutor) addReturnValue(command string, returnValue []string) {
	executor.returnValues[command] = returnValue
}

/*
	Asserts that a given slice of commands have been called executed against the command executor.
 */
func (executor *MockedCommandExecutor) assertCommands(expectedCommands []string) {
	if !reflect.DeepEqual(executor.executedCommands, expectedCommands) {
		executor.t.Errorf("Executed command were '%v', expected '%v'", executor.executedCommands, expectedCommands)
	}
}

func createTestClient(t *testing.T) (*memClient, *MockedCommandExecutor) {
	executor := &MockedCommandExecutor{t, []string{}, map[string][]string{}, false}
	client := &memClient{
		server:   "foo",
		executor: executor,
	}
	return client, executor
}

// Actual tests

func TestMemClient(t *testing.T) {
	_, err := MemClient("foo:1234")
	if err == nil {
		t.Errorf("Memclient should return an error for foo:1234")
	}
}

func TestGet(t *testing.T) {
	client, executor := createTestClient(t)
	client.Get("testkey")
	executor.assertCommands([]string{"get testkey\r\n"})
}

func TestSet(t *testing.T) {
	client, executor := createTestClient(t)
	client.Set("testkey", "testval", 123)
	executor.assertCommands([]string{"set testkey 0 123 7\r\ntestval\r\n"})
}

func TestVersion(t *testing.T) {
	client, executor := createTestClient(t)
	version := client.Version()
	executor.assertCommands([]string{"version \r\n"})
	if version != "UNKNOWN" {
		t.Errorf("Received version does not match expected version (%v!=%v)", version, "VERSION myversion.1234")
	}

	executor.addReturnValue("version \r\n", []string{"VERSION myversion.1234"})
	version = client.Version()
	if version != "VERSION myversion.1234" {
		t.Errorf("Received version does not match expected version (%v!=%v)", version, "VERSION myversion.1234")
	}
}

func TestStats(t *testing.T) {
	client, executor := createTestClient(t)
	// return some random stats
	returnStats := []string{"STAT time 1446586044", "STAT version 1.4.14 (Ubuntu)", "STAT libevent 2.0.21-stable"}
	executor.addReturnValue("stats\r\n", returnStats)

	stats := client.Stats()

	// validate that the result is correct and that the expected commands were executed
	expectedStats := []Stat{
		{"time", "1446586044"},
		{"version", "1.4.14 (Ubuntu)"},
		{"libevent", "2.0.21-stable"},
	}

	if !reflect.DeepEqual(stats, expectedStats) {
		t.Errorf("Returned cache stats incorrect (%v!=%v)", stats, expectedStats)
	}

	executor.assertCommands([]string{"stats\r\n"})
}

func TestListKeys(t *testing.T) {
	// setup testcase
	client, executor := createTestClient(t)
	executor.addReturnValue("stats items\r\n", []string{"STAT items:1:number 4"})

	executor.addReturnValue("stats cachedump 1 4\n", []string{
		"ITEM location [9 b; 1539093795 s]",
		"ITEM profession [9 b; 1539088675 s]",
		"ITEM age [4 b; 1539088575 s]",
		"ITEM username [8 b; 1539088375 s]",
	})

	keys := client.ListKeys()

	// validate that the result is correct and that the expected commands were executed
	expectedKeys := []Key{
		{Original: "ITEM location [9 b; 1539093795 s]", Name: "location", Expiry: 1539093795},
		{Original: "ITEM profession [9 b; 1539088675 s]", Name: "profession", Expiry: 1539088675},
		{Original: "ITEM age [4 b; 1539088575 s]", Name: "age", Expiry: 1539088575},
		{Original: "ITEM username [8 b; 1539088375 s]", Name: "username", Expiry: 1539088375},
	}

	if (!reflect.DeepEqual(keys, expectedKeys)) {
		fmt.Println(keys)
		t.Errorf("Returned cache keys incorrect (%v!=%v)", keys, expectedKeys)
	}

	executor.assertCommands([]string{"stats items\r\n", "stats cachedump 1 4\n"})
}

func TestStat(t *testing.T) {
	// setup testcase
	client, executor := createTestClient(t)
	returnStats := []string{"STAT time 1446586044", "STAT version 1.4.14 (Ubuntu)", "STAT libevent 2.0.21-stable"}
	executor.addReturnValue("stats\r\n", returnStats)

	time, ok := client.Stat("time")

	expectedTime := Stat{"time", "1446586044"}
	if !ok || !reflect.DeepEqual(time, expectedTime) {
		t.Errorf("Returned cache stat incorrect (%v!=%v)", time, expectedTime)
	}
	executor.assertCommands([]string{"stats\r\n"})
}
