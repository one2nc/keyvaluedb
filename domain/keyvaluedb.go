package domain

import (
	"fmt"
	"strconv"
)

type KeyValueDB struct {
	storage             map[string]interface{}
	isMultiBlockStarted bool
	cmds                []Command
}

func NewKeyValueDB() *KeyValueDB {
	return &KeyValueDB{storage: make(map[string]interface{})}
}

func (kvdb *KeyValueDB) Execute(cmd Command) interface{} {
	_, err := cmd.Validate()
	if err != nil {
		return err.Error()
	}

	if kvdb.isMultiBlockStarted && !cmd.isTerminatorCmd() {
		kvdb.enqueue(cmd)
		return "QUEUED"
	}

	switch cmd.Name {
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
		var outputs []interface{}
		for key, val := range kvdb.storage {
			outputs = append(outputs, fmt.Sprintf("SET %v %v", key, val))
		}
		return outputs
	case SET:
		kvdb.storage[cmd.Key] = cmd.Value
		return "OK"
	case GET:
		v, ok := kvdb.storage[cmd.Key]
		if !ok {
			return nil
		}
		return v
	case DEL:
		_, ok := kvdb.storage[cmd.Key]
		if !ok {
			return 0
		}
		delete(kvdb.storage, cmd.Key)
		return 1
	case INCR:
		v, ok := kvdb.storage[cmd.Key]
		if !ok {
			newResult := "1"
			kvdb.storage[cmd.Key] = newResult
			return newResult
		}

		currentValue, err := strconv.Atoi(v.(string))
		if err != nil {
			return "(error) ERR value is not an integer or out of range"
		}

		incrementedValue := fmt.Sprintf("%v", currentValue+1)
		kvdb.storage[cmd.Key] = incrementedValue
		return incrementedValue
	case INCRBY:
		v, ok := kvdb.storage[cmd.Key]
		if !ok {
			newResult := cmd.Value
			kvdb.storage[cmd.Key] = newResult
			return newResult
		}

		currentValue, err := strconv.Atoi(v.(string))
		if err != nil {
			return "(error) ERR value is not an integer or out of range"
		}

		resultValue, err := strconv.Atoi(cmd.Value.(string))
		if err != nil {
			return "(error) ERR value is not an integer or out of range"
		}

		incrementedValue := fmt.Sprintf("%v", currentValue+resultValue)

		kvdb.storage[cmd.Key] = incrementedValue
		return incrementedValue
	}

	return fmt.Errorf("(error) ERR unknown command '%s'", cmd.Key)
}

func (kvdb *KeyValueDB) enqueue(cmd Command) {
	kvdb.cmds = append(kvdb.cmds, cmd)
}

func (kvdb *KeyValueDB) executeCommands() interface{} {
	var outputs []interface{}
	for _, cmd := range kvdb.cmds {
		outputs = append(outputs, kvdb.Execute(cmd))
	}
	return outputs
}
