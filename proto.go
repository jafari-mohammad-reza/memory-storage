package main

import (
	"errors"
	"strings"
)

type Command interface {
	Execute() error
}
type SetCommand struct {
	key string
	val string
}

func (c *SetCommand) Execute() error {
	key, val := c.key, c.val
	if key == "" {
		return errors.New("must enter a key")
	}
	if val == "" {
		return errors.New("must enter a value")
	}
	return nil
}

type GetCommand struct {
	key string
}

func (c *GetCommand) Execute() error {
	key := c.key
	if key == "" {
		return errors.New("must enter a key")
	}
	return nil
}

type DelCommand struct {
	key string
}

func (c *DelCommand) Execute() error {
	key := c.key
	if key == "" {
		return errors.New("must enter a key")
	}
	return nil
}

func parseCommand(cmd string) (Command, error) {
	inp := strings.Fields(cmd)
	switch strings.ToLower(strings.TrimSpace(inp[0])) {
	case "set":
		if len(inp) != 3 {
			return nil, errors.New("set command needs both key and value")
		}
		return &SetCommand{
			key: inp[1],
			val: inp[2],
		}, nil
	case "get":
		if len(inp) != 2 {
			return nil, errors.New("get command needs key")
		}
		return &GetCommand{
			key: inp[1],
		}, nil
	case "del":
		if len(inp) != 2 {
			return nil, errors.New("delete command needs key")
		}
		return &DelCommand{
			key: inp[1],
		}, nil
	default:
		return nil, errors.New("invalid command")
	}
}
