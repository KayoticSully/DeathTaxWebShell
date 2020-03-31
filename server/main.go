package main

import (
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/KayoticSully/DeathTaxWebShell/server/deathtax"
	"github.com/fasthttp/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

var sessions = []*deathtax.Session{}
var instanceFactory *deathtax.PooledFactory

const webDir = "/usr/local/share/deathtax/web"

func main() {
	assetDir := path.Join(webDir, "assets")

	http.HandleFunc("/", index)
	http.Handle("/assets/", http.StripPrefix(strings.TrimRight("/assets/", "/"), http.FileServer(http.Dir(assetDir))))
	http.HandleFunc("/api", api)

	instanceFactory = deathtax.NewPooledFactory(3, 10)

	log.Println("Server Up!")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(webDir, "index.html"))
}

func api(w http.ResponseWriter, r *http.Request) {
	// TODO: Websocket connection handling (errors, close event, etc)
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer wsConn.Close()

	session := instanceFactory.GetInstance()
	sessions = append(sessions, session)

	session.RunWebsocketProxy(wsConn)
}
