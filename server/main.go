package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/KayoticSully/DeathTaxWebShell/server/deathtax"
	"github.com/fasthttp/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func main() {

	http.HandleFunc("/", index)
	http.Handle("/assets/", http.StripPrefix(strings.TrimRight("/assets/", "/"), http.FileServer(http.Dir("../site/assets"))))
	http.HandleFunc("/api", api)

	log.Println("Server Up!")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../site/index.html")
}

func api(w http.ResponseWriter, r *http.Request) {

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer wsConn.Close()

	dt := deathtax.New()
	dtResp := <-dt.Output

	err = wsConn.WriteMessage(websocket.TextMessage, []byte(dtResp))
	if err != nil {
		log.Println("write:", err)
		return
	}

	for {
		msgType, msg, err := wsConn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		dt.Input <- string(msg)
		log.Println("Waiting for output")
		dtResp = <-dt.Output

		log.Printf("Output: %s\n", msg)

		err = wsConn.WriteMessage(msgType, []byte(dtResp))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
