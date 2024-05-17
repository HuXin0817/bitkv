package api

import (
	"context"
	"github.com/HuXin0817/bitkv/interval/errors"
	"github.com/HuXin0817/bitkv/interval/pb"
	"github.com/zeromicro/go-zero/zrpc"
)

type Client struct {
	client pb.ServeClient
}

func NewClient(url string) *Client {
	rpcConf := zrpc.RpcClientConf{Endpoints: []string{url}}
	return &Client{
		client: pb.NewServeClient(zrpc.MustNewClient(rpcConf).Conn()),
	}
}

func (c *Client) Put(k, v string) (err error) {
	if err = errors.CheckKey(k); err != nil {
		return err
	}
	if err = errors.CheckValue(v); err != nil {
		return err
	}
	if _, err = c.client.Put(context.Background(), &pb.Pair{Key: k, Value: v}); err != nil {
		return err
	}
	return nil
}

func (c *Client) Delete(k string) (err error) {
	return c.Put(k, "")
}

func (c *Client) Get(k string) (v string, err error) {
	if err = errors.CheckKey(k); err != nil {
		return "", err
	}
	r, err := c.client.Get(context.Background(), &pb.Key{Key: k})
	if err != nil {
		return "", err
	}
	return r.Value, nil
}
