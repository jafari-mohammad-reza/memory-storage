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
	if s.log {
		defer s.Log("KEYS")
	}
	return s.listStorage.ShowAll()
}
func (s *MemoryStorage) RGet(value []byte) ([]byte, error) {
	if s.log {
		defer s.Log("KEYS")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listStorage.Get(value)
}
func (s *MemoryStorage) RSet(value []byte) error {
	if s.log {
		defer s.Log("KEYS")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listStorage.Insert(value)
}
func (s *MemoryStorage) RDel(value []byte) error {
	if s.log {
		defer s.Log("KEYS")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listStorage.Delete(value)
}
func (s *MemoryStorage) Keys() ([]string, error) {
	if s.log {
		defer s.Log("KEYS")
	}
	var keys []string
	for key, _ := range s.mapStorage {
		keys = append(keys, key)
	}
	return keys, nil
}

func (s *MemoryStorage) Get(key string) (*[]byte, error) {
	if s.log {
		defer s.Log("GET", key)
	}
	val := s.mapStorage[key]
	if val == nil {
		return nil, fmt.Errorf("there is no item with key of %s", key)
	}
	return &val, nil
}

func (s *MemoryStorage) Del(key string) error {
	if s.log {
		defer s.Log("DEL", key)
	}
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
	logFile, err := os.OpenFile(fmt.Sprintf("storage_%d.log", s.Id), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	data := LogData{
		Action: command,
		Args:   make([]string, 0),
	}
	for _, a := range arg {
		data.Args = append(data.Args, a.(string))
	}
	log, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = logFile.Write(log)
	if err != nil {
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
}

type StorageCore interface {
	GetStorage(id int, log bool) Storage
	RecoverFromLogs(storageId int) error
}

type MemoryStorageCore struct {
	Storages []*MemoryStorage
}

func NewMemoryStorageCore() *MemoryStorageCore {
	storages := make([]*MemoryStorage, 0)
	storages = append(storages, NewStorage(0, false))
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

func (sc *MemoryStorageCore) RecoverFromLogs(storageId int) error {
	// this method will recover given id storage from logs of that storage id and append it to core
	return nil
}
