package cluster

import (
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
	"my_redis/lib/logger"
	"my_redis/resp/client"
)

type connectionFactory struct {
	peer string //地址
}

func (f connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	c, err := client.MakeClient(f.peer)
	if err != nil{
		logger.Error(err)
		return nil, err
	}

	c.Start()

	return pool.NewPooledObject(c), nil
}

func (f connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	c, ok := object.Object.(*client.Client)
	if ok == false{
		return errors.New("type mismatch")
	}

	c.Close()

	return nil
}

func (f connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (f connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (f connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

