package main

import (
	"keyvaluedb/domain"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGetAndDelKeyValue(t *testing.T) {
	// SET key1 value1
	// SET key2 value2
	// GET key1 // should return value1
	// DEL key2 // deletes key2 and its value
	// GET key2 // should return nil
	// SET key2 newvalue2

	keyValueDB := domain.NewKeyValueDB()

	assert.Equal(t, "value1", keyValueDB.Execute(domain.NewCommand(domain.SET, "key1", "value1")))
	assert.Equal(t, "value2", keyValueDB.Execute(domain.NewCommand(domain.SET, "key2", "value2")))
	assert.Equal(t, "value1", keyValueDB.Execute(domain.NewCommand(domain.GET, "key1")))
	assert.Equal(t, "value2", keyValueDB.Execute(domain.NewCommand(domain.DEL, "key2")))
	assert.Equal(t, "", keyValueDB.Execute(domain.NewCommand(domain.GET, "key2")))
	assert.Equal(t, "newvalue2", keyValueDB.Execute(domain.NewCommand(domain.SET, "key2", "newvalue2")))
}

func TestBasicOperations(t *testing.T) {
	// SET counter 0
	// INCR counter // increments a "counter" key, if present and returns incremented value, in this case: 1
	// GET counter // returns 1
	// INCRBY counter 10 // increment by 10, returns 11
	// INCR foo // automatically creates a new key called "foo" and increments it by 1 and thus returns 1

	keyValueDB := domain.NewKeyValueDB()

	assert.Equal(t, "0", keyValueDB.Execute(domain.NewCommand(domain.SET, "counter", 0)))
	intVar, _ := strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.INCR, "counter")))
	assert.Equal(t, 1, intVar)
	intVar, _ = strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.GET, "counter")))
	assert.Equal(t, 1, intVar)
	intVar, _ = strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.INCRBY, "counter", 10)))
	assert.Equal(t, 11, intVar)
	intVar, _ = strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.INCR, "foo")))
	assert.Equal(t, 1, intVar)
}

func TestMultipleCmdsCase1Test(t *testing.T) {
	// MULTI // starts a multi line commands
	// INCR foo // queues this command, doesn't execute it immediately
	// SET key1 value1 // queues this command, doesn't execute it immediately
	// EXEC // execute all queued commands and returns output of all commands in an array, thus returns: [1 value1]

	keyValueDB := domain.NewKeyValueDB()

	assert.Equal(t, "OK", keyValueDB.Execute(domain.NewCommand(domain.MULTI)))
	assert.Equal(t, "OK", keyValueDB.Execute(domain.NewCommand(domain.INCR, "foo")))
	assert.Equal(t, "OK", keyValueDB.Execute(domain.NewCommand(domain.SET, "key1", "value1")))
	assert.Equal(t, "[1, value1]", keyValueDB.Execute(domain.NewCommand(domain.EXEC)))
}

func TestMultipleCmdsCase2(t *testing.T) {
	// MULTI // starts a multi line commands
	// INCR foo // queues this command, doesn't execute it immediately
	// SET key1 value1 // queues this command, doesn't execute it immediately
	// DISCARD // discard all queued commands
	// GET key1 // returns nil as key1 doesn't exists

	keyValueDB := domain.NewKeyValueDB()

	assert.Equal(t, "OK", keyValueDB.Execute(domain.NewCommand(domain.MULTI)))
	assert.Equal(t, "OK", keyValueDB.Execute(domain.NewCommand(domain.INCR, "foo")))
	assert.Equal(t, "OK", keyValueDB.Execute(domain.NewCommand(domain.SET, "key1", "value1")))
	assert.Equal(t, "OK", keyValueDB.Execute(domain.NewCommand(domain.DISCARD)))
	assert.Equal(t, "", keyValueDB.Execute(domain.NewCommand(domain.GET, "key1")))
}

func TestCompactCase1(t *testing.T) {
	// SET counter 10
	// INCR counter
	// INCR counter
	// SET foo bar
	// GET counter // returns 12
	// INCR counter
	// COMPACT // this should return following output
	// SET counter 13
	// SET foo bar

	keyValueDB := domain.NewKeyValueDB()

	assert.Equal(t, "10", keyValueDB.Execute(domain.NewCommand(domain.SET, "counter", 10)))
	intVar, _ := strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.INCR, "counter")))
	assert.Equal(t, 11, intVar)
	intVar, _ = strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.INCR, "counter")))
	assert.Equal(t, 12, intVar)
	assert.Equal(t, "bar", keyValueDB.Execute(domain.NewCommand(domain.SET, "foo", "bar")))
	intVar, _ = strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.GET, "counter")))
	assert.Equal(t, 12, intVar)
	intVar, _ = strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.INCR, "counter")))
	assert.Equal(t, 13, intVar)
	assert.Equal(t, "OK", keyValueDB.Execute(domain.NewCommand(domain.COMPACT)))
}

func TestCompactCase2(t *testing.T) {
	// INCR counter // returns 1
	// INCRBY counter 10 // returns 11
	// GET counter // returns 11
	// DEL counter // deletes counter
	// COMPACT // this should compact to empty output as there's no keys present in the DB

	keyValueDB := domain.NewKeyValueDB()

	intVar, _ := strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.INCR, "counter")))
	assert.Equal(t, 1, intVar)
	intVar, _ = strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.INCRBY, "counter", 10)))
	assert.Equal(t, 11, intVar)
	intVar, _ = strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.GET, "counter")))
	assert.Equal(t, 11, intVar)
	intVar, _ = strconv.Atoi(keyValueDB.Execute(domain.NewCommand(domain.DEL, "counter")))
	assert.Equal(t, 11, intVar)
	assert.Equal(t, "OK", keyValueDB.Execute(domain.NewCommand(domain.COMPACT)))
}
