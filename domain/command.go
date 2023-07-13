package domain

import (
	"fmt"
)

const (
	SET     string = "SET"
	GET     string = "GET"
	DEL     string = "DEL"
	INCR    string = "INCR"
	INCRBY  string = "INCRBY"
	MULTI   string = "MULTI"
	EXEC    string = "EXEC"
	DISCARD string = "DISCARD"
	COMPACT string = "COMPACT"
)

type command struct {
	Name  string
	Key   string
	Value interface{}
}

func NewCommand(name string, args ...interface{}) command {
	var key string
	var value interface{}
	if len(args) > 1 {
		value = args[1]
	}
	if len(args) > 0 {
		key = fmt.Sprintf("%v", args[0])
	}

	return command{
		Name:  name,
		Key:   key,
		Value: value,
	}
}

func (c command) isTerminatorCmd() bool {
	switch c.Name {
	case EXEC, DISCARD:
		return true
	}
	return false
}

func (c command) Validate() (bool, error) {
	switch c.Name {
	case SET, INCRBY:
		cmd := "set"
		if c.Name == INCRBY {
			cmd = "incrby"
		}
		if c.Value == nil {
			return false, fmt.Errorf("(error) ERR wrong number of arguments for '%s' command", cmd)
		}
		return true, nil
	case GET, DEL, INCR:
		cmd := "get"
		switch c.Name {
		case DEL:
			cmd = "del"
		case INCR:
			cmd = "incrby"
		}
		if c.Key == "" {
			return false, fmt.Errorf("(error) ERR wrong number of arguments for '%s' command", cmd)
		}
		return true, nil
	case MULTI, EXEC, DISCARD, COMPACT:
		return true, nil
	}

	return false, fmt.Errorf("(error) ERR unknown command `%s`, with args beginning with: '%s', '%s',", c.Name, c.Key, c.Value)
}
