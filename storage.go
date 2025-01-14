package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type LinkedListNode struct {
	value    []byte
	NextNode *LinkedListNode
	PrevNode *LinkedListNode
}

type LinkedList struct {
	head *LinkedListNode
}

func (l *LinkedList) Insert(value []byte) error {
	newNode := &LinkedListNode{
		value:    value,
		NextNode: nil,
		PrevNode: nil,
	}
	if l.head == nil {
		l.head = newNode
		return nil
	}
	current := l.head
	for current.NextNode != nil {
		current = current.NextNode
	}
	current.NextNode = newNode
	newNode.PrevNode = current
	return nil
}

func (l *LinkedList) Delete(value []byte) error {
	if l.head == nil {
		return errors.New("list is empty")
	}

	current := l.head
	for current != nil && !bytes.Equal(current.value, value) {
		current = current.NextNode
	}

	if current == nil {
		return errors.New("value not found in the list")
	}

	if current.PrevNode != nil {
		current.PrevNode.NextNode = current.NextNode
	} else {
		l.head = current.NextNode
	}

	if current.NextNode != nil {
		current.NextNode.PrevNode = current.PrevNode
	}

	return nil
}

func (l *LinkedList) ShowAll() ([]byte, error) {
	if l.head == nil {
		return nil, errors.New("list is empty")
	}

	var buffer bytes.Buffer
	current := l.head
	for current != nil {
		buffer.Write(current.value)
		if current.NextNode != nil {
			buffer.Write([]byte(", "))
		}
		current = current.NextNode
	}

	return buffer.Bytes(), nil
}
func (l *LinkedList) Get(value []byte) ([]byte, error) {
	if l.head == nil {
		return nil, errors.New("list is empty")
	}

	current := l.head
	for current != nil {
		if bytes.Equal(current.value, value) {
			return current.value, nil
		}
		current = current.NextNode
	}

	return nil, errors.New("value not found in the list")
}

type MemoryStorage struct {
	Id          int
	mapStorage  map[string][]byte
	listStorage *LinkedList
	log         bool
	mu          sync.RWMutex
}

func NewStorage(id int, log bool) *MemoryStorage {
	return &MemoryStorage{
		Id:          id,
		mapStorage:  make(map[string][]byte),
		listStorage: &LinkedList{},
		log:         log,
	}
}

func (s *MemoryStorage) RKeys() ([]byte, error) {
	return s.listStorage.ShowAll()
}
func (s *MemoryStorage) RGet(value []byte) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listStorage.Get(value)
}
func (s *MemoryStorage) RSet(value []byte) error {
	if s.log {
		defer s.Log("RSET")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listStorage.Insert(value)
}
func (s *MemoryStorage) RDel(value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listStorage.Delete(value)
}
func (s *MemoryStorage) Keys() ([]string, error) {
	var keys []string
	for key, _ := range s.mapStorage {
		keys = append(keys, key)
	}
	return keys, nil
}

func (s *MemoryStorage) Get(key string) (*[]byte, error) {
	val := s.mapStorage[key]
	if val == nil {
		return nil, fmt.Errorf("there is no item with key of %s", key)
	}
	return &val, nil
}

func (s *MemoryStorage) Del(key string) error {
	val := s.mapStorage[key]
	if val == nil {
		return fmt.Errorf("there is no item with key of %s", key)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.mapStorage, key)
	return nil
}

func (s *MemoryStorage) Set(key string, value []byte) error {
	if s.log {
		defer s.Log("SET", key, value)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.mapStorage[key] = value
	return nil
}

type LogData struct {
	Action string
	Args   []string
}

func (s *MemoryStorage) Log(command string, arg ...interface{}) error {
	logFile, err := os.OpenFile(fmt.Sprintf("storage_%d.json", s.Id), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	stat, err := logFile.Stat()
	if err != nil {
		return err
	}
	if stat.Size() == 0 {
		if _, err := logFile.Write([]byte("[]")); err != nil {
			return err
		}
	}

	var logs []LogData
	if _, err := logFile.Seek(0, 0); err != nil {
		return err
	}
	if err := json.NewDecoder(logFile).Decode(&logs); err != nil {
		return err
	}

	data := LogData{
		Action: command,
		Args:   make([]string, 0),
	}
	for _, a := range arg {
		switch v := a.(type) {
		case string:
			data.Args = append(data.Args, v)
		case []uint8:
			data.Args = append(data.Args, string(v))
		default:
			continue
		}
	}

	logs = append(logs, data)

	if err := logFile.Truncate(0); err != nil {
		return err
	}
	if _, err := logFile.Seek(0, 0); err != nil {
		return err
	}
	if err := json.NewEncoder(logFile).Encode(logs); err != nil {
		return err
	}

	return nil
}

type Storage interface {
	Keys() ([]string, error)
	Get(key string) (*[]byte, error)
	Del(key string) error
	Set(key string, value []byte) error
	RSet(value []byte) error
	RKeys() ([]byte, error)
	RGet(value []byte) ([]byte, error)
	RDel(value []byte) error
	Log(command string, arg ...interface{}) error
	RecoverFromLogs() error
}

type StorageCore interface {
	GetStorage(id int, log bool) Storage
}

type MemoryStorageCore struct {
	Storages []*MemoryStorage
}

func NewMemoryStorageCore() *MemoryStorageCore {
	storages := make([]*MemoryStorage, 0)
	storages = append(storages, NewStorage(0, true))
	return &MemoryStorageCore{
		Storages: storages,
	}
}

func (sc *MemoryStorageCore) GetStorage(id int, log bool) Storage {
	var storage *MemoryStorage
	for _, st := range sc.Storages {
		if st.Id == id {
			storage = st
			break
		}
	}
	if storage == nil {
		newSt := NewStorage(id, log)
		sc.Storages = append(sc.Storages, newSt)
		return newSt
	}
	return storage
}

func (s *MemoryStorage) RecoverFromLogs() error {
	fmt.Println("recovering from logs")
	logData, err := os.ReadFile(fmt.Sprintf("storage_%d.json", s.Id))
	if err != nil {
		return err
	}
	var logEntries []LogData
	err = json.Unmarshal(logData, &logEntries)
	if err != nil {
		return err
	}
	for _, entry := range logEntries {
		fmt.Printf("Processing Action: %s, Args: %v\n", entry.Action, entry.Args)
		switch entry.Action {
		case "SET":
			if len(entry.Args) >= 2 {
				s.Set(entry.Args[0], []byte(entry.Args[1]))
			}
		case "RSET":
			if len(entry.Args) >= 1 {
				s.RSet([]byte(entry.Args[1]))
			}
		default:
			fmt.Printf("Unknown action: %s\n", entry.Action)
		}
	}
	return nil
}
