package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type StorageValue struct {
	Value     string
	CreatedAt time.Time
}

type MemoryStorage struct {
	Id         int
	memoryData map[string]StorageValue
	log        bool
}

func NewStorage(id int, log bool) *MemoryStorage {
	return &MemoryStorage{
		Id:         id,
		memoryData: make(map[string]StorageValue),
		log:        log,
	}
}

func (s *MemoryStorage) Keys() ([]string, error) {
	if s.log {
		defer s.Log("KEYS")
	}
	var keys []string
	for key, _ := range s.memoryData {
		keys = append(keys, key)
	}
	return keys, nil
}

func (s *MemoryStorage) Get(key string) (*StorageValue, error) {
	if s.log {
		defer s.Log("GET", key)
	}
	val := s.memoryData[key]
	if val.Value == "" {
		return nil, fmt.Errorf("there is no item with key of %s", key)
	}
	return &val, nil
}

func (s *MemoryStorage) Del(key string) error {
	if s.log {
		defer s.Log("DEL", key)
	}
	val := s.memoryData[key]
	if val.Value == "" {
		return fmt.Errorf("there is no item with key of %s", key)
	}
	delete(s.memoryData, key)
	return nil
}

func (s *MemoryStorage) Set(key string, value string) error {
	if s.log {
		defer s.Log("SET", key, value)
	}
	s.memoryData[key] = StorageValue{
		Value:     value,
		CreatedAt: time.Now().UTC(),
	}
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
	Get(key string) (*StorageValue, error)
	Del(key string) error
	Set(key string, value string) error
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
