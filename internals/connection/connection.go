package connection

import (
	"errors"

	"github.com/00000vish/filebeam/internals/config"
)

type Connection interface {
	IsConnected() bool
	IsHosting() bool
	Create(config *config.Config)
	Host() error
	Connect() error
	Send(data []byte) error
	Recieve(callback func(data []byte))
}

func Create(config *config.Config) (Connection, error) {
	var socket Connection = nil

	switch config.Connection {
	case "ws":
		socket = &Websocket{}
		break
	}

	if socket == nil {
		return nil, errors.New("Invalid connection type.")
	}

	socket.Create(config)

	return socket, nil
}
