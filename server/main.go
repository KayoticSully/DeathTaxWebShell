package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/KayoticSully/DeathTaxWebShell/server/deathtax"
	"github.com/fasthttp/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

var sessions = []*deathtax.Session{}
var instanceFactory *deathtax.PooledFactory
var gitShas *GitShas

const publicAPIPort = ":5000"
const webDir = "/usr/local/share/deathtax/web"

func main() {
	// Start this as early as possible
	instanceFactory = deathtax.NewPooledFactory(5)

	// Boot up k8s internal endpoints
	go startHealthCheckAPI()

	assetDir := path.Join(webDir, "assets")
	publicMux := http.NewServeMux()

	publicMux.HandleFunc("/", index)
	publicMux.Handle("/assets/", http.StripPrefix(strings.TrimRight("/assets/", "/"), http.FileServer(http.Dir(assetDir))))
	publicMux.HandleFunc("/api", api)

	gitShas = loadGitShas(fmt.Sprintf("%s/%s", webDir, "_meta"))

	log.Println("Server Up!")
	log.Fatal(http.ListenAndServe(publicAPIPort, publicMux))
}

func index(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join(webDir, "index.html")
	tmpl, _ := template.ParseFiles(fp)
	tmpl.ExecuteTemplate(w, "index", gitShas)
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
