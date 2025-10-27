package internal

import (
	"context"
	kv "api-gateway/internal/proto/kv/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	kv kv.KeyValueClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{kv: kv.NewKeyValueClient(conn)}, nil
}

func (c *Client) Put(ctx context.Context, key string, value []byte) error {
	_, err := c.kv.Put(ctx, &kv.PutRequest{Key: key, Value: value})
	return err
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := c.kv.Get(ctx, &kv.GetRequest{Key: key})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.kv.Delete(ctx, &kv.DeleteRequest{Key: key})
	return err
}
