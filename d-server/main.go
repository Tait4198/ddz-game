package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// todo tap联想

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func urlParam(key string, r *http.Request) (string, error) {
	params, ok := r.URL.Query()[key]
	if !ok || len(params[0]) < 1 {
		return "", fmt.Errorf("url param '%s' is missing", key)
	}
	param := params[0]
	return param, nil
}

func main() {
	center := newCenter()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		usr, err1 := urlParam("usr", r)
		pwd, err2 := urlParam("pwd", r)

		if err1 != nil || err2 != nil {
			w.WriteHeader(403)
			log.Println("usr或pwd为空")
			return
		}

		clientId := newClientId(usr, pwd)
		if center.checkClientId(clientId) {
			w.WriteHeader(401)
			return
		}

		conn, err := upgrade.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		newClient(usr, clientId, center, conn)
	})
	_ = http.ListenAndServe(":8080", nil)
}
