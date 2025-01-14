package main

import (
	"context"
	"testing"
)

var server *Server

func TestMain(m *testing.M) {
	server = NewServer(Config{
		ServerListenAddr: "8001",
	})
	server.Start()
	defer server.Stop()

	m.Run()
}

func newTestClient() *Client {
	return NewClient(":8001")
}

func setKey(t *testing.T, client *Client, key, value string) {
	t.Helper()
	_, err := client.Set(context.Background(), key, value)
	if err != nil {
		t.Fatalf("failed to set key %q: %v", key, err)
	}
}

func TestSetCommand(t *testing.T) {
	client := newTestClient()
	response, err := client.Set(context.Background(), "hello", "world")
	if err != nil {
		t.Fatal(err)
	}
	expectedResp := `"SET executed key: hello , value: world"`
	if response != expectedResp {
		t.Errorf("unexpected response: got %q, want %q", response, expectedResp)
	}
}

func TestGetCommand(t *testing.T) {
	client := newTestClient()
	setKey(t, client, "hello", "world")

	response, err := client.Get(context.Background(), "hello")
	if err != nil {
		t.Fatal(err)
	}
	expectedResp := `"world"`
	if response != expectedResp {
		t.Errorf("unexpected response: got %q, want %q", response, expectedResp)
	}
}

func TestGetCommandCases(t *testing.T) {
	client := newTestClient()
	setKey(t, client, "hello", "world")

	tests := []struct {
		name         string
		key          string
		expectedResp string
		expectError  bool
	}{
		{"KeyExists", "hello", `"world"`, false},
		{"KeyNotFound", "nonexistent", `there is no item with key of nonexistent`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := client.Get(context.Background(), tt.key)
			if (err != nil) != tt.expectError {
				t.Fatalf("unexpected error: %v", err)
			}
			if response != tt.expectedResp {
				t.Errorf("unexpected response: got %q, want %q", response, tt.expectedResp)
			}
		})
	}
}

func TestDelCommand(t *testing.T) {
	client := newTestClient()
	setKey(t, client, "hello", "world")

	response, err := client.Del(context.Background(), "hello")
	if err != nil {
		t.Fatal(err)
	}
	expectedResp := `"DEL executed key: hello"`
	if response != expectedResp {
		t.Errorf("unexpected response: got %q, want %q", response, expectedResp)
	}
}

func TestKeysCommand(t *testing.T) {
	client := newTestClient()
	setKey(t, client, "hello", "world")
	setKey(t, client, "world", "hello")

	response, err := client.Keys(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expectedResp := `["hello","world"]`
	if response != expectedResp {
		t.Errorf("unexpected response: got %q, want %q", response, expectedResp)
	}
}
