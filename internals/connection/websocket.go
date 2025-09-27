package connection

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/00000vish/filebeam/internals/config"
	"github.com/00000vish/filebeam/internals/logger"

	"github.com/gorilla/websocket"
)

const WS_PATH = "data"

type Websocket struct {
	config               *config.Config
	createTimeOutSeconds int64
	isConnected          bool
	isHosting            bool
	connection           *websocket.Conn
	subscribers          []func(data []byte)
}

func (w *Websocket) IsHosting() bool {
	return w.isHosting
}

func (w *Websocket) IsConnected() bool {
	return w.isConnected
}

func (w *Websocket) Create(config *config.Config) {
	w.config = config
	w.createTimeOutSeconds = 2
}

func (w *Websocket) tryRunTillError(run func() error, errorCallback func(error)) (bool, error) {
	errored := make(chan error)

	go func() {
		err := run()
		errored <- err
	}()

	time.Sleep(time.Duration(w.createTimeOutSeconds) * time.Second)

	select {
	case res := <-errored:
		if res == nil {
			return true, nil
		}
		return false, res
	case <-time.After(time.Duration(w.createTimeOutSeconds) * time.Second):
		go func() {
			err := <-errored
			if err != nil {
				errorCallback(err)
			}
		}()
	}

	return true, nil
}

func (w *Websocket) Host() error {

	if w.config == nil {
		return errors.New("call create first")
	}

	if w.isHosting {
		return nil
	}

	logger.Log("Creating websocket server.")

	var action = func() error {
		http.HandleFunc(fmt.Sprintf("/%s", WS_PATH), w.dataHandler)
		err := http.ListenAndServe(fmt.Sprintf(":%s", w.config.Port), nil)
		return err
	}

	var callback = func(err error) {
		logger.Log(err.Error())
		w.isHosting = false
	}

	started, err := w.tryRunTillError(action, callback)

	w.isHosting = started

	return err
}

func (w *Websocket) Connect() error {

	if w.config == nil {
		return errors.New("call create first")
	}

	if w.isConnected {
		return nil
	}

	logger.Log("Tring to connect websocket.")

	var action = func() error {
		url := fmt.Sprintf("ws://%s:%s/%s", w.config.IpAddress, w.config.Port, WS_PATH)

		var err error = nil

		w.connection, _, err = websocket.DefaultDialer.Dial(url, nil)

		if err != nil {
			return err
		}

		return err
	}

	var callback = func(err error) {
		logger.Log(err.Error())
		w.isConnected = false
	}

	started, err := w.tryRunTillError(action, callback)

	w.isConnected = started

	return err
}

func (w *Websocket) Send(data []byte) error {
	err := w.connection.WriteMessage(1, data)
	if err != nil {
		logger.Log(err.Error())
	}
	return nil
}

func (w *Websocket) Recieve(callback func(data []byte)) {
	w.subscribers = append(w.subscribers, callback)
}

var upgrader = websocket.Upgrader{}

func (w *Websocket) dataHandler(writer http.ResponseWriter, reader *http.Request) {
	c, err := upgrader.Upgrade(writer, reader, nil)
	if err != nil {
		exit(err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			exit(err)
			return
		}
		for _, callback := range w.subscribers {
			go callback(message)
		}
	}
}

func exit(err error) {
	logger.Log(err.Error())
	logger.Log("Exiting now...")
	os.Exit(0)
}
