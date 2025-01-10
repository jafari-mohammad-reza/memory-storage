package main

import "context"

type Client struct {
	connString string
}

func NewClient(connString string) *Client {
	return &Client{
		connString,
	}
}


func (c *Client) Keys(ctx context.Context) error {
	return nil
}

func (c *Client) Set(ctx context.Context, key string, value string) error {
	return nil
}

func (c *Client) Get(ctx context.Context, key string) error {
	return nil
}

func (c *Client) Del(ctx context.Context, key string) error {
	return nil
}
