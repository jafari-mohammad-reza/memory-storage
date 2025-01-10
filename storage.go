package main

import (
	"time"
)

type StorageValue struct {
	Value     string
	CreatedAt time.Time
}

type Storage struct {
	Id         int
	memoryData map[string]StorageValue
}

func NewStorage(id int) *Storage {
	return &Storage{
		Id:         id,
		memoryData: make(map[string]StorageValue),
	}
}

type StorageCore struct {
	Storages []*Storage
}

func NewStorageCore() *StorageCore {
	storages := make([]*Storage, 0)
	storages = append(storages, NewStorage(0))
	return &StorageCore{
		Storages: storages,
	}
}

func (sc *StorageCore) GetStorage(id int) (*Storage, error) {
	var storage *Storage
	for _, st := range sc.Storages {
		if st.Id == id {
			storage = st
			break
		}
	}
	if storage == nil {
		newSt := NewStorage(id)
		sc.Storages = append(sc.Storages, newSt)
		return newSt, nil
	}
	return storage, nil
}
func (sc *StorageCore) RecoverFromLogs(storageId int) error {
	// this method will recover given id storage from logs of that storage id and append it to core
	return nil
}
