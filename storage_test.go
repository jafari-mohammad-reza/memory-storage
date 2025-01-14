package main

import (
	"bytes"
	"testing"
)

func TestMemoryStorage_SetAndGet(t *testing.T) {
	storage := NewStorage(1, false)
	err := storage.Set("key1", []byte("value1"))
	if err != nil {
		t.Fatalf("unexpected error on Set: %v", err)
	}
	value, err := storage.Get("key1")
	if err != nil {
		t.Fatalf("unexpected error on Get: %v", err)
	}
	if !bytes.Equal(*value, []byte("value1")) {
		t.Errorf("expected value 'value1', got '%s'", *value)
	}
}

func TestMemoryStorage_Del(t *testing.T) {
	storage := NewStorage(1, false)
	err := storage.Set("key1", []byte("value1"))
	if err != nil {
		t.Fatalf("unexpected error on Set: %v", err)
	}
	err = storage.Del("key1")
	if err != nil {
		t.Fatalf("unexpected error on Del: %v", err)
	}
	_, err = storage.Get("key1")
	if err == nil {
		t.Errorf("expected error on Get after delete, got nil")
	}
}

func TestLinkedList_InsertAndShowAll(t *testing.T) {
	list := &LinkedList{}
	err := list.Insert([]byte("value1"))
	if err != nil {
		t.Fatalf("unexpected error on Insert: %v", err)
	}
	err = list.Insert([]byte("value2"))
	if err != nil {
		t.Fatalf("unexpected error on Insert: %v", err)
	}
	values, err := list.ShowAll()
	if err != nil {
		t.Fatalf("unexpected error on ShowAll: %v", err)
	}
	expected := []byte("value1, value2")
	if !bytes.Equal(values, expected) {
		t.Errorf("expected values '%s', got '%s'", expected, values)
	}
}

func TestLinkedList_Delete(t *testing.T) {
	list := &LinkedList{}
	list.Insert([]byte("value1"))
	list.Insert([]byte("value2"))
	err := list.Delete([]byte("value1"))
	if err != nil {
		t.Fatalf("unexpected error on Delete: %v", err)
	}
	values, err := list.ShowAll()
	if err != nil {
		t.Fatalf("unexpected error on ShowAll: %v", err)
	}
	expected := []byte("value2")
	if !bytes.Equal(values, expected) {
		t.Errorf("expected values '%s', got '%s'", expected, values)
	}
}
