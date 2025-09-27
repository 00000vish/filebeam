package endpoint

import (
	"bytes"
	"fmt"

	"github.com/00000vish/filebeam/internals/connection"
	"github.com/00000vish/filebeam/internals/logger"
)

type Endpoint struct {
	connection connection.Connection
	onConnect  func()
	onReceive  func([]byte)
}

func (e *Endpoint) Create(connection connection.Connection) error {
	e.connection = connection

	connection.Recieve(e.receiver)

	err := connection.Host()
	if err != nil {
		return err
	}

	return nil
}

func (e *Endpoint) Connect(onConnect func(), onRecieve func([]byte)) {
	e.onConnect = onConnect
	e.onReceive = onRecieve

	e.tryConnect()
}

func (e *Endpoint) tryConnect() {
	if e.connection.IsConnected() {
		return
	}

	err := e.connection.Connect()
	if err == nil {
		e.connection.Send(HAND_SHAKE[:])

		logger.Log("Connected to remote filebeam service.")

		e.onConnect()
	}
}

func (e *Endpoint) receiver(data []byte) {
	if bytes.Equal(data[0:PROTOCOL_SIZE], HAND_SHAKE[:]) {
		e.tryConnect()
		return
	}

	if bytes.Equal(data[0:PROTOCOL_SIZE], DATA_SEND[:]) {
		e.onReceive(data[PROTOCOL_SIZE:])
		return
	}

	logger.Log(fmt.Sprintf("Recieved unknown data : %s", string(data[:])))
}

func (e *Endpoint) Send(data []byte) {
	e.sendData(DATA_SEND, data)
}

func (e *Endpoint) sendData(protocol []byte, data []byte) {
	toSend := append(protocol, data...)

	err := e.connection.Send(toSend)
	if err != nil {
		logger.Log(err.Error())
	}
}
