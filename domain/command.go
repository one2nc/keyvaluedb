package domain

import (
	"fmt"
)

type cmd int

const (
	SET cmd = iota
	GET
	DEL
	INCR
	INCRBY
	MULTI
	EXEC
	DISCARD
	COMPACT
)

type command struct {
	c    cmd
	args []interface{}
}

func NewCommand(c cmd, args ...interface{}) command {
	return command{
		c:    c,
		args: args,
	}
}

func (c command) Run() (*entity, error) {

	switch c.c {
	case SET, INCRBY:
		return &entity{
			Key:   fmt.Sprintf("%v", c.args[0]),
			Value: fmt.Sprintf("%v", c.args[1]),
		}, nil
	case GET, DEL, INCR:
		return &entity{
			Key: fmt.Sprintf("%v", c.args[0]),
		}, nil
	case MULTI, EXEC, DISCARD, COMPACT:
		return nil, nil
	}

	return nil, fmt.Errorf("command not yet implemented")
}
