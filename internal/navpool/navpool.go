package navpool

import (
	"errors"
	"log"
)

type PoolApi struct {
	client *PoolClient
}

var (
	ErrPoolConnection = errors.New("Could not connect to the NavPool API")
)

func NewPoolApi(host string, network string) (*PoolApi, error) {
	poolClient, err := NewClient(host, network)
	if err != nil {
		log.Print(err)
		return nil, ErrPoolConnection
	}

	return &PoolApi{poolClient}, nil
}
