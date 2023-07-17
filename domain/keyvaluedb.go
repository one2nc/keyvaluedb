package domain

import (
	"fmt"
	"keyvaluedb/storage"
	"strconv"
)

type KeyValueDB struct {
	storage             storage.Storage
	isMultiBlockStarted bool
	cmds                []Command
}

func NewKeyValueDB(storage storage.Storage) KeyValueDB {
	return KeyValueDB{storage: storage}
}

func (kvdb *KeyValueDB) Execute(dbIndex int, cmd Command) (int, interface{}) {
	_, err := cmd.Validate()
	if err != nil {
		return dbIndex, err.Error()
	}

	if kvdb.isMultiBlockStarted && !cmd.isTerminatorCmd() {
		kvdb.enqueue(cmd)
		return dbIndex, "QUEUED"
	}

	switch cmd.Name {
	case SELECT:
		dbIndex, err = kvdb.storage.Select(cmd.Key)
		if err != nil {
			return dbIndex, err.Error()
		}
		return dbIndex, "OK"
	case MULTI:
		kvdb.isMultiBlockStarted = true
		return dbIndex, "OK"
	case DISCARD:
		kvdb.isMultiBlockStarted = false
		kvdb.cmds = nil
		return dbIndex, "OK"
	case EXEC:
		kvdb.isMultiBlockStarted = false
		return dbIndex, kvdb.executeCommands(dbIndex)
	case COMPACT:
		var outputs []interface{}
		for keyVal := range kvdb.storage.GetAll(dbIndex) {
			outputs = append(outputs, fmt.Sprintf("SET %s", keyVal))
		}
		return dbIndex, outputs
	case SET:
		kvdb.storage.Set(dbIndex, cmd.Key, cmd.Value)
		return dbIndex, "OK"
	case GET:
		return dbIndex, kvdb.storage.Get(dbIndex, cmd.Key)
	case DEL:
		return dbIndex, kvdb.storage.Del(dbIndex, cmd.Key)
	case INCR:
		v := kvdb.storage.Get(dbIndex, cmd.Key)
		if v == nil {
			newResult := "1"
			kvdb.storage.Set(dbIndex, cmd.Key, newResult)
			return dbIndex, newResult
		}

		currentValue, err := strconv.Atoi(v.(string))
		if err != nil {
			return dbIndex, "(error) ERR value is not an integer or out of range"
		}

		incrementedValue := fmt.Sprintf("%v", currentValue+1)
		kvdb.storage.Set(dbIndex, cmd.Key, incrementedValue)
		return dbIndex, incrementedValue
	case INCRBY:
		v := kvdb.storage.Get(dbIndex, cmd.Key)
		if v == nil {
			newResult := cmd.Value
			kvdb.storage.Set(dbIndex, cmd.Key, newResult)
			return dbIndex, newResult
		}

		currentValue, err := strconv.Atoi(v.(string))
		if err != nil {
			return dbIndex, "(error) ERR value is not an integer or out of range"
		}

		resultValue, err := strconv.Atoi(cmd.Value.(string))
		if err != nil {
			return dbIndex, "(error) ERR value is not an integer or out of range"
		}

		incrementedValue := fmt.Sprintf("%v", currentValue+resultValue)

		kvdb.storage.Set(dbIndex, cmd.Key, incrementedValue)
		return dbIndex, incrementedValue
	}

	return dbIndex, fmt.Errorf("(error) ERR unknown command '%s'", cmd.Key)
}

func (kvdb *KeyValueDB) enqueue(cmd Command) {
	kvdb.cmds = append(kvdb.cmds, cmd)
}

func (kvdb *KeyValueDB) executeCommands(dbIndex int) interface{} {
	var outputs []interface{}
	for _, cmd := range kvdb.cmds {
		_, result := kvdb.Execute(dbIndex, cmd)
		outputs = append(outputs, result)
	}
	kvdb.cmds = nil
	return outputs
}
