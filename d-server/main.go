package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

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
	port := 8080
	flag.IntVar(&port, "port", 8080, "服务器端口")
	flag.Parse()

	log.Printf("使用[%d]端口启动服务", port)

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
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err == nil {
		log.Println(err)
	}
}
