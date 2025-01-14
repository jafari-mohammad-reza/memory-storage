package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	connString string
	conn       net.Conn
	mu         sync.Mutex
}

func NewClient(connString string) *Client {
	return &Client{connString: connString}
}

func (c *Client) ensureConnection() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		conn, err := net.Dial("tcp", c.connString)
		if err != nil {
			return fmt.Errorf("failed to connect to server: %w", err)
		}
		c.conn = conn
	}
	return nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

func (c *Client) sendCommand(ctx context.Context, command string) (string, error) {

	if err := c.ensureConnection(); err != nil {
		return "", err
	}

	c.mu.Lock()
	_, err := c.conn.Write([]byte(command + "\r\n"))
	c.mu.Unlock()
	if err != nil {
		c.Close()
		return "", fmt.Errorf("failed to send command: %w", err)
	}

	return c.readResponse(ctx)
}

func (c *Client) readResponse(ctx context.Context) (string, error) {
	buffer := make([]byte, 1024)
	c.mu.Lock()
	defer c.mu.Unlock()

	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := c.conn.Read(buffer)
	if err != nil {
		c.Close()
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	return string(buffer[:n]), nil
}

func (c *Client) Set(ctx context.Context, key string, value string) (string, error) {
	command := fmt.Sprintf("SET %s %s", key, value)
	return c.sendCommand(ctx, command)
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	command := fmt.Sprintf("GET %s", key)
	return c.sendCommand(ctx, command)
}

func (c *Client) Del(ctx context.Context, key string) (string, error) {
	command := fmt.Sprintf("DEL %s", key)
	return c.sendCommand(ctx, command)
}
func (c *Client) Keys(ctx context.Context) (string, error) {
	command := fmt.Sprintf("KEYS")
	return c.sendCommand(ctx, command)
}
