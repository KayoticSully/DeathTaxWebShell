package deathtax

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/fasthttp/websocket"
)

// Session powershell process object
type Session struct {
	process *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
}

// NewSession returns a new running DeathTax process
func NewSession() *Session {
	s := &Session{
		// TODO: Refactor to standard / configurable path
		process: exec.Command("pwsh", "/usr/local/share/deathtax/DeathTax"),
	}

	return s
}

// Run starts the session and proxies input/output to a websocket
func (s *Session) Run() {
	var err error

	// Input and Output pipes need to be created before the go
	// routines start. Otherwise data will be missed between process
	// start and go routine start.  There is an unknowable delay between
	// telling a go routine to start and when it actually starts.
	s.stdin, err = s.process.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	s.stdout, err = s.process.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Start data pumps. Now that the pipes have been created, these
	// can get started whenever the scheduled pleases.

	log.Println("Starting SubProcess")

	if err := s.process.Start(); err != nil {
		log.Fatal(err)
	}
}

// RunWebsocketProxy pipes input and output channels from/to the websocket
func (s *Session) RunWebsocketProxy(wsConn *websocket.Conn) {
	go inputPump(s.stdin, wsConn)
	outputPump(s.stdout, wsConn)
}

func inputPump(stdin io.WriteCloser, wsConn *websocket.Conn) {
	// stdin is created outside of the go routine but it should
	// be setup to close when the go routine exits
	defer stdin.Close()

	var err error
	var msg []byte

	for {
		_, msg, err = wsConn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		_, err = io.WriteString(stdin, string(msg))
	}
}

func outputPump(stdout io.ReadCloser, wsConn *websocket.Conn) {
	// stdout is created outside of the go routine but it should
	// be setup to close when the go routine exits
	defer stdout.Close()

	var text []byte
	var err error
	stdScanner := bufio.NewScanner(stdout)

	for stdScanner.Scan() {
		text = []byte(stdScanner.Text())
		if strings.TrimSpace(string(text)) != "" {
			err = wsConn.WriteMessage(websocket.TextMessage, text)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}