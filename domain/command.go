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
	SELECT  string = "SELECT"
)

type Command struct {
	Name  string
	Key   string
	Value interface{}
}

func NewCommand(name string, args ...interface{}) Command {
	var key string
	var value interface{}
	if len(args) > 1 {
		value = args[1]
	}
	if len(args) > 0 {
		key = fmt.Sprintf("%v", args[0])
	}

	return Command{
		Name:  name,
		Key:   key,
		Value: value,
	}
}

func (c Command) isTerminatorCmd() bool {
	switch c.Name {
	case EXEC, DISCARD:
		return true
	}
	return false
}

func (c Command) Validate() (bool, error) {
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
	case SELECT:
		cmd := "select"
		if c.Key == "" {
			return false, fmt.Errorf("(error) ERR wrong number of arguments for '%s' command", cmd)
		}
		return true, nil
	case MULTI, EXEC, DISCARD, COMPACT:
		return true, nil
	}

	params := ""
	if c.Key != "" {
		params = fmt.Sprintf("`%s`,", c.Key)
		if c.Value != nil {
			params = fmt.Sprintf("`%s`, `%v`,", c.Key, c.Value)
		}
	}

	return false, fmt.Errorf("(error) ERR unknown command `%s`, with args beginning with: %s", c.Name, params)
}
