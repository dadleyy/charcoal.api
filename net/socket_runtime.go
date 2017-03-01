package net

import "io"
import "fmt"
import "bytes"
import "net/http"
import "github.com/gorilla/websocket"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/activity"

type SocketRuntime struct {
	*log.Logger
	Messages chan activity.Message
}

func (r *SocketRuntime) checkOrigin(_ *http.Request) bool {
	return true
}

func (r *SocketRuntime) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: r.checkOrigin, ReadBufferSize: 1024, WriteBufferSize: 1024}
	connection, err := upgrader.Upgrade(res, req, nil)

	if res.Header().Get("status") != "" {
		return
	}

	if err != nil {
		r.Debugf("unable to connect websocket: %s", err.Error())
		return
	}

	r.Debugf("websocket open - %v", connection.RemoteAddr())
	res.Header().Set("status", "SENT")

	connection.WriteMessage(websocket.TextMessage, []byte("hello world"))

	open := true

	for open {
		select {
		case message := <-r.Messages:
			ping := bytes.NewBuffer([]byte(fmt.Sprintf("received: %s", message.Verb)))
			r.Debugf("received message: %s", message.Verb)

			w, err := connection.NextWriter(websocket.TextMessage)

			if err != nil {
				r.Debugf("connection appears closed: %s", err.Error())
				open = false
				break
			}

			if _, err := io.Copy(w, ping); err != nil {
				r.Debugf("connection appears closed: %s", err.Error())
				open = false
				break
			}

			if err := w.Close(); err != nil {
				r.Debugf("connection appears closed: %s", err.Error())
				open = false
				break
			}

			r.Debugf("successfully wrote to %v", connection.RemoteAddr())
		}
	}

	r.Debugf("closing connection w/ %s", connection.RemoteAddr())
	connection.Close()
}
