package main

import (
	"errors"
	"fmt"
	"strings"
)

type Command interface {
	Execute(Storage) (interface{}, error)
}

type KeysCommand struct{}

func (c *KeysCommand) Execute(s Storage) (interface{}, error) {
	fmt.Println("Keys command got executed")
	keys, err := s.Keys()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

type SetCommand struct {
	key string
	val string
}

func (c *SetCommand) Execute(s Storage) (interface{}, error) {
	key, val := c.key, c.val
	if key == "" {
		return nil, errors.New("must enter a key")
	}
	if val == "" {
		return nil, errors.New("must enter a value")
	}

	err := s.Set(key, []byte(val))
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("SET executed key: %s , value: %s", c.key, c.val), nil
}

type GetCommand struct {
	key string
}

func (c *GetCommand) Execute(s Storage) (interface{}, error) {
	key := c.key
	if key == "" {
		return nil, errors.New("must enter a key")
	}
	val, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	fmt.Printf("GET executed key: %s", c.key)
	return string(*val), nil
}

type DelCommand struct {
	key string
}

func (c *DelCommand) Execute(s Storage) (interface{}, error) {
	key := c.key
	if key == "" {
		return nil, errors.New("must enter a key")
	}
	err := s.Del(key)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("DEL executed key: %s", c.key), nil
}

type RKeysCommand struct {
	key string
}

func (c *RKeysCommand) Execute(s Storage) (interface{}, error) {
	fmt.Println("RKeys command got executed")
	keys, err := s.RKeys()
	if err != nil {
		return nil, err
	}
	return string(keys), nil
}

type RSetCommand struct {
	val string
}

func (c *RSetCommand) Execute(s Storage) (interface{}, error) {
	val := c.val

	if val == "" {
		return nil, errors.New("must enter a value")
	}

	err := s.RSet([]byte(val))
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("RSET executed value: %s", val), nil
}

type RDelCommand struct {
	val string
}

func (c *RDelCommand) Execute(s Storage) (interface{}, error) {
	val := c.val
	if val == "" {
		return nil, errors.New("must enter a val")
	}
	err := s.RDel([]byte(val))
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("DEL executed value: %s", val), nil
}

type RGetCommand struct {
	val string
}

func (c *RGetCommand) Execute(s Storage) (interface{}, error) {
	val := c.val
	if val == "" {
		return nil, errors.New("must enter a val")
	}
	foundVal, err := s.RGet([]byte(val))
	if err != nil {
		return nil, err
	}
	return string(foundVal), nil
}

type RecoverCommand struct {
}

func (c *RecoverCommand) Execute(s Storage) (interface{}, error) {
	err := s.RecoverFromLogs()
	if err != nil {
		return nil, err
	}
	return "recovered successfully", nil
}

func parseCommand(cmd string) (Command, error) {
	inp := strings.Fields(cmd)
	if len(inp) != 0 {
		switch strings.ToLower(strings.TrimSpace(inp[0])) {
		case "keys":
			return &KeysCommand{}, nil
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
		case "rkey":
			return &RKeysCommand{}, nil
		case "rget":
			if len(inp) != 2 {
				return nil, errors.New("RGet command needs value")
			}
			return &RGetCommand{
				val: inp[1],
			}, nil
		case "rdel":
			if len(inp) != 2 {
				return nil, errors.New("RDel command needs value")
			}
			return &RDelCommand{
				val: inp[1],
			}, nil
		case "rset":
			if len(inp) != 2 {
				return nil, errors.New("RSet command needs value")
			}
			return &RSetCommand{
				val: inp[1],
			}, nil
		case "recover":
			return &RecoverCommand{}, nil
		default:
			return nil, errors.New("invalid command")
		}
	}
	return nil, errors.New("empty command.\n")
}
