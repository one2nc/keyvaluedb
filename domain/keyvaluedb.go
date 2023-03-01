package domain

import (
	"fmt"
	"strconv"
	"strings"
)

type keyValueDB struct {
	storage             map[string]entity
	isMultiBlockStarted bool
	cmds                []command
}

func NewKeyValueDB() *keyValueDB {
	return &keyValueDB{storage: make(map[string]entity)}
}

func (kvdb *keyValueDB) Execute(cmd command) string {

	switch cmd.c {
	case MULTI:
		kvdb.isMultiBlockStarted = true
		return "OK"
	case DISCARD:
		kvdb.isMultiBlockStarted = false
		kvdb.cmds = nil
		return "OK"
	case EXEC:
		kvdb.isMultiBlockStarted = false
		return kvdb.executeCommands()
	case COMPACT:
		for _, val := range kvdb.storage {
			fmt.Println("SET", val.Key, val.Value)
		}
		return "OK"
	}

	if kvdb.isMultiBlockStarted {
		kvdb.enqueue(cmd)
		return "OK"
	}

	result, err := cmd.Run()
	if err != nil {
		return "ERROR"
	}

	switch cmd.c {
	case SET:
		kvdb.storage[result.Key] = *result
		return result.Value
	case GET:
		v, ok := kvdb.storage[result.Key]
		if !ok {
			return ""
		}
		return v.Value
	case DEL:
		v, ok := kvdb.storage[result.Key]
		if !ok {
			return ""
		}
		delete(kvdb.storage, result.Key)
		return v.Value
	case INCR:
		v, ok := kvdb.storage[result.Key]
		if !ok {
			newResult := entity{
				Key:   result.Key,
				Value: "1",
			}
			kvdb.storage[result.Key] = newResult
			return newResult.Value
		}

		currentValue, err := strconv.Atoi(v.Value)
		if err != nil {
			return err.Error()
		}

		incrementedValue := fmt.Sprintf("%v", currentValue+1)
		newResult := entity{
			Key:   result.Key,
			Value: incrementedValue,
		}
		kvdb.storage[result.Key] = newResult
		return incrementedValue
	case INCRBY:
		v, ok := kvdb.storage[result.Key]
		if !ok {
			newResult := entity{
				Key:   result.Key,
				Value: "1",
			}
			kvdb.storage[result.Key] = newResult
			return newResult.Value
		}

		currentValue, err := strconv.Atoi(v.Value)
		if err != nil {
			return err.Error()
		}

		resultValue, err := strconv.Atoi(result.Value)
		if err != nil {
			return err.Error()
		}

		incrementedValue := fmt.Sprintf("%v", currentValue+resultValue)

		newResult := entity{
			Key:   result.Key,
			Value: incrementedValue,
		}
		kvdb.storage[result.Key] = newResult
		return incrementedValue
	}

	return ""
}

func (kvdb *keyValueDB) enqueue(cmd command) {
	kvdb.cmds = append(kvdb.cmds, cmd)
}

func (kvdb *keyValueDB) executeCommands() string {
	var outputs []string
	for _, cmd := range kvdb.cmds {
		outputs = append(outputs, kvdb.Execute(cmd))
	}
	return fmt.Sprintf("[%s]", strings.Join(outputs, ", "))
}
